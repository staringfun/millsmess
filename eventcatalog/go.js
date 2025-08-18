const fsPromises = require('fs/promises')
const $RefParser = require('@apidevtools/json-schema-ref-parser');
const gofmt = require('gofmt.js');
const utils = require('@eventcatalog/sdk');
const { getEvents, getCommands } = utils.default(__dirname);

function patternToCharList(pattern) {
    // Extract character class inside [...]
    const match = pattern.match(/\[([^\]]+)\]/);
    if (!match) return [];

    const classPart = match[1];
    const chars = [];

    for (let i = 0; i < classPart.length; i++) {
        const char = classPart[i];

        // Handle ranges like A-Z
        if (i + 2 < classPart.length && classPart[i + 1] === '-') {
            const start = classPart.charCodeAt(i);
            const end = classPart.charCodeAt(i + 2);
            for (let c = start; c <= end; c++) {
                chars.push(String.fromCharCode(c));
            }
            i += 2; // Skip the range
        } else {
            chars.push(char);
        }
    }

    return chars;
}

function formatName(name) {
    return capitalizeFirst(name).replace(/Id/g, "ID");
}

function formatEnumName(name) {
    return name
        .split(/[:._]/)
        .map(part => part.charAt(0).toUpperCase() + part.slice(1))
        .join('');
}

function formatType(property) {
    if (property.title) return property.title
    if (property.goType) return property.goType
    switch (property.type) {
        case 'boolean':
            return 'bool'
        case 'string':
            if (property.format === 'date' || property.format === 'date-time') {
                return 'time.Time'
            }
            break
        case 'array':
            return `[]${formatType(property.items)}`
        case 'object':
            if (property.additionalProperties) {
                if (!property.properties) {
                    const key = property.goKeyType || 'any'
                    return `map[${key}]${formatType(property.additionalProperties)}`
                }
            } else {
                return 'any'
            }
        break
    }
    return property.type
}

const capitalizeFirst = (str) => str.charAt(0).toUpperCase() + str.slice(1)

class FileWriter {
    constructor(schemas) {
        this.data = []
        this.schemas = schemas
    }

    appendLine(line) {
        this.data.push(line)
    }

    getText() {
        return this.data.join('\n')
    }

    appendLinesWithTabs(number, ...lines) {
        lines.forEach(l => {
            this.appendLine(`${new Array(number).fill('\t').join('')}${l}`)
        })
    }

    writeBaseType(schema) {
        const schemaType = schema.goType || schema.type
        const ref = schema.title[0].charAt(0).toLowerCase()
        if (schema.goEqualsTo) {
            this.appendLine(`type ${schema.title} = ${schema.goEqualsTo}`)
        } else {
            this.appendLine(`type ${schema.title} ${schemaType}`)
        }

        this.appendLine(`func (${ref} ${schema.title}) ${capitalizeFirst(schemaType)}() ${schemaType} {`)
        this.appendLinesWithTabs(1,`return ${schemaType}(${ref})`)
        this.appendLine(`}`)

        const emptyValue = schema.goEmptyValue ?? ''
        this.appendLine(`func (${ref} ${schema.title}) IsEmpty() bool {`)
        this.appendLinesWithTabs(1, `return ${ref} == ${JSON.stringify(emptyValue)}`)
        this.appendLine(`}`)

        if (schema.enum) {
            this.appendLine('const (')
            schema.enum.forEach(v => {
                this.appendLinesWithTabs(1, `${schema.title}${formatEnumName(v)} ${schema.title} = ${JSON.stringify(v)}`)
            })
            this.appendLine(')')
            this.appendLine(`var All${schema.title} = []${schema.title}{${schema.enum.map(e => `${schema.title}${formatEnumName(e)}`).join(', ')}}`)
            this.appendLine(`func (${ref} ${schema.title}) IsValid() bool {`)
            this.appendLinesWithTabs(1, `for _, ${ref}${ref} := range All${schema.title} {`)
            this.appendLinesWithTabs(2, `if ${ref}${ref} == ${ref} {`)
            this.appendLinesWithTabs(3, `return true`)
            this.appendLinesWithTabs(2, `}`)
            this.appendLinesWithTabs(1, `}`)
            this.appendLinesWithTabs(1, `return false`)
            this.appendLine(`}`)
            return
        }

        if (schema.pattern && !schema.goNoRunes) {
            const chars = patternToCharList(schema.pattern)
            if (chars.length) {
                this.appendLine(`var ${schema.title}Runes = []rune{`)
                chars.forEach(c => {
                    this.appendLinesWithTabs(1, `'${c}',`)
                })
                this.appendLine('}')
            }
        }

        if (schema.minLength) {
            this.appendLine(`const ${schema.title}MinLength = ${schema.minLength}`)
        }

        if (schema.maxLength) {
            this.appendLine(`const ${schema.title}MaxLength = ${schema.maxLength}`)
        }

        this.appendLine(`func (${ref} ${schema.title}) IsValid() bool {`)
        if (schema.minLength) {
            this.appendLinesWithTabs(1, `if len(${ref}) < ${schema.title}MinLength {`)
            this.appendLinesWithTabs(2, `return false`)
            this.appendLinesWithTabs(1, `}`)
        }
        if (schema.maxLength) {
            this.appendLinesWithTabs(1, `if len(${ref}) > ${schema.title}MaxLength {`)
            this.appendLinesWithTabs(2, `return false`)
            this.appendLinesWithTabs(1, `}`)
        }
        if (schema.pattern && !schema.goNoRunes) {
            this.appendLine(`outer:`)
            this.appendLinesWithTabs(1, `for _, ${ref}${ref} := range []rune(${ref}) {`)
            this.appendLinesWithTabs(2, `for _, ${ref}${ref}${ref} := range ${schema.title}Runes {`)
            this.appendLinesWithTabs(3, `if ${ref}${ref} == ${ref}${ref}${ref} {`)
            this.appendLinesWithTabs(4, `continue outer`)
            this.appendLinesWithTabs(3, `}`)
            this.appendLinesWithTabs(2, `}`)
            this.appendLinesWithTabs(1,
                `return false`,
                '}'
            )
        }
        this.appendLinesWithTabs(1, `return ${ref} != ${JSON.stringify(emptyValue)}`)
        this.appendLine(`}`)
    }

    writeStruct(schema, name) {
        const ref = (name || schema.title).charAt(0).toLowerCase()
        this.appendLine(`type ${name || schema.title} struct {`)
        Object.keys(schema.properties).forEach((p) => {
            const property = schema.properties[p]
            if (schema.goImplements) {
                for (let i = 0; i < this.schemas.length; i++) {
                    const e = this.schemas[i]
                    if (e.title === schema.goImplements) {
                        const keys = Object.keys(e.properties)
                        if (keys.length) {
                            for (let j = 0; j < keys.length; j++) {
                                const ep = keys[j]
                                if (ep === p) {
                                    return
                                }
                            }
                        } else {
                          break
                        }
                    }
                }
            }
            const required = schema.required?.length ? property.goInterface || schema.required.some(r => r === p) : property.goInterface
            this.appendLinesWithTabs(1, `${formatName(p)} ${required ? '' : '*'}${formatType(property)} \`json:"${p},omitempty"\``)
        })
        if (schema.goImplements) {
            this.appendLinesWithTabs(1, `*Base${schema.goImplements}`)
        }
        this.appendLine(`}`)
        if (!schema.goNoIsValid) {
            this.appendLine(`func (${ref} *${name || schema.title}) IsValid() bool {`)
            this.appendLinesWithTabs(1, `if ${ref} == nil {`)
            this.appendLinesWithTabs(2, `return false`)
            this.appendLinesWithTabs(1, `}`)
            Object.keys(schema.properties).forEach((p) => {
                const property = schema.properties[p]
                if (property.format === 'date' || property.format === 'date-time') return
                if (property.type === 'boolean') return
                if (property.goInterface) return;
                if (property.type === 'object') {
                    if (property.additionalProperties) {
                        this.appendLinesWithTabs(1, `for _, ${ref}${ref} := range ${ref}.${formatName(p)} {`)
                        this.appendLinesWithTabs(2, `if !${ref}${ref}.IsValid() {`)
                        this.appendLinesWithTabs(3, `return false`)
                        this.appendLinesWithTabs(2, `}`)
                        this.appendLinesWithTabs(1, `}`)
                        return;
                    } else {
                        if (!property.properties) {
                            return;
                        }
                    }
                }
                if (property.type === 'array') {
                    this.appendLinesWithTabs(1, `for _, ${ref}${ref} := range ${ref}.${formatName(p)} {`)
                    this.appendLinesWithTabs(2, `if !${ref}${ref}.IsValid() {`)
                    this.appendLinesWithTabs(3, `return false`)
                    this.appendLinesWithTabs(2, `}`)
                    this.appendLinesWithTabs(1, `}`)
                    return;
                }
                const required = schema.required?.length ? property.goInterface || schema.required.some(r => r === p) : property.goInterface
                this.appendLinesWithTabs(1, `if ${required ? '' : `${ref}.${formatName(p)} != nil && `}!${ref}.${formatName(p)}.IsValid() {`)
                this.appendLinesWithTabs(2, `return false`)
                this.appendLinesWithTabs(1, `}`)
            })
            this.appendLinesWithTabs(1, `return true`)
            this.appendLine(`}`)
        }
    }

    writeInterface(schema) {
        this.appendLine(`type ${schema.title} interface {`)
        Object.keys(schema.properties).forEach((p) => {
            const property = schema.properties[p]
            const required = schema.required?.length ? schema.required.some(r => r === p) : false
            this.appendLinesWithTabs(1, `Get${formatName(p)}() ${required ? '' : '*'}${formatType(property)}`)
            this.appendLinesWithTabs(1, `Set${formatName(p)}(${required ? '' : '*'}${formatType(property)})`)
        })
        this.appendLine(`}`)
        if (!schema.noGoBase) {
            const ref = schema.title.charAt(0).toLowerCase()
            this.writeStruct({...schema, goNoIsValid: true}, `Base${schema.title}`)
            Object.keys(schema.properties).forEach((p) => {
                const property = schema.properties[p]
                const required = schema.required?.length ? schema.required.some(r => r === p) : false
                this.appendLine(`func (${ref} *Base${schema.title}) Get${formatName(p)}() ${required ? '' : '*'}${formatType(property)} {`)
                this.appendLinesWithTabs(1, `if ${ref} == nil {`)
                this.appendLinesWithTabs(2, `var ${ref}${ref} ${required ? '' : '*'}${formatType(property)}`)
                this.appendLinesWithTabs(2, `return ${ref}${ref}`)
                this.appendLinesWithTabs(1, `}`)
                this.appendLinesWithTabs(1, `return ${ref}.${formatName(p)}`)
                this.appendLine('}')
                this.appendLine(`func (${ref} *Base${schema.title}) Set${formatName(p)}(${ref}${ref} ${required ? '' : '*'}${formatType(property)}) {`)
                this.appendLinesWithTabs(1, `if ${ref} == nil {`)
                this.appendLinesWithTabs(2, `return`)
                this.appendLinesWithTabs(1, `}`)
                this.appendLinesWithTabs(1, `${ref}.${formatName(p)} = ${ref}${ref}`)
                this.appendLine('}')
            })
        }
    }

    writeSchema(schema) {
        this.appendLine('')
        this.appendLine('')
        if (!schema.properties) {
            this.writeBaseType(schema)
        } else {
            if (schema.goInterface) {
                this.writeInterface(schema)
            } else {
                this.writeStruct(schema)
            }
        }
    }

    writeSchemas() {
        this.schemas.forEach(s => {
            this.writeSchema(s)
        })
    }
}
// The config value is your EventCatalog config
// The options value is the configuration you pass to your plugin
module.exports = async (config, options = {}) => {
    const {
        packageName = 'types',
        outPath = 'generated.go',
        schemaPath = 'schema.json',
        preamble = [],
    } = options;

    const [schema, commands, events] = await Promise.all([
        $RefParser.dereference(schemaPath),
        getCommands(),
        getEvents(),
    ])
    const defs = Object.values(schema['$defs'])


    const writer = new FileWriter(defs)
    preamble.forEach(p => {
        writer.appendLine(p)
    })
    writer.appendLine(`// Package ${packageName}: Code generated; DO NOT EDIT;`)
    writer.appendLine(`package ${packageName}`)
    writer.appendLine('')
    writer.appendLine(`import "time"`)

    writer.writeSchemas()

    writer.appendLine('')
    writer.appendLine('')
    writer.writeBaseType({
        title: 'SocketMessageTypeCommand',
        type: 'string',
        enum: commands.filter(c => c.channels.some(s => s.id === 'WS')).map(c => `v${c.version}:${c.name}`)
    })
    writer.writeBaseType({
        title: 'SocketMessageTypeEvent',
        type: 'string',
        enum: events.filter(c => c.channels.some(s => s.id === 'WS')).map(c => `v${c.version}:${c.name}`)
    })

    const text = writer.getText()
    const formatted = gofmt(text)
    await fsPromises.writeFile(outPath, formatted ? formatted : `// FILE CONTAINS ERRORS!\n${text}`)
}