const fsPromises = require('fs/promises')
const $RefParser = require('@apidevtools/json-schema-ref-parser');
const init = require('@wasm-fmt/gofmt').default;
const { format } = require('@wasm-fmt/gofmt');
const utils = require('@eventcatalog/sdk');
const { getEvents, getCommands } = utils.default(__dirname);
const { FileWriter } = require('./generator')

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
        init()
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
        enum: Array.from(new Set(commands.filter(c => c.channels.some(cc => cc.id === 'WS')).map(c => `v${c.version}:${c.name}`)))
    })
    writer.writeBaseType({
        title: 'SocketMessageTypeEvent',
        type: 'string',
        enum: Array.from(new Set(events.filter(e => e.channels.some(cc =>  cc.id === 'WS')).map(e => `v${e.version}:${e.name}`)))
    })
    writer.writeBaseType({
        title: 'TopicName',
        type: 'string',
        enum: Array.from(new Set([
            ...(commands.filter(c => c.channels.some(cc =>  cc.id !== 'WS')).map(c => c.name)),
            ...(events.filter(e => e.channels.some(cc =>  cc.id !== 'WS')).map(e => e.name)),
        ]))
    })

    const text = writer.getText()
    await fsPromises.writeFile(outPath, format(text))
}