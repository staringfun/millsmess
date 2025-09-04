const fsPromises = require('fs/promises')
const init = require('@wasm-fmt/gofmt').default;
const {format} = require('@wasm-fmt/gofmt');
const utils = require('@eventcatalog/sdk');
const { getEvents, getCommands } = utils.default(__dirname);
const { FileWriter, formatEnumName } = require('./generator')

function writeMessage(writer, message, channels, opts) {
    const {
        refName,
        structName,
        engineName,
        typesPackage,
        marshallerName,
    } = opts
    const schema = require(`./${message.schemaPath}`)
    message.channels?.forEach(c => {
        if (!channels[c.id]) {
            channels[c.id] = {events:[]}
        }
        channels[c.id].events.push({
            id: message.id,
            version: message.version,
            schemaTitle: schema.title,
        })
    })
    writer.appendLine(``)
    writer.appendLine(``)
    writer.appendLine(`func (${refName} *${structName}) PublishV${message.version}${message.id}(data ${typesPackage}.${schema.title}, attributes map[string]string, config PublishConfig, ctx context.Context) error {`)
    writer.appendLinesWithTabs(1, `bytes, err := ${refName}.${marshallerName}.Marshal(data)`)
    writer.appendLinesWithTabs(1, `if err != nil {`)
    writer.appendLinesWithTabs(2, `return err`)
    writer.appendLinesWithTabs(1, `}`)
    writer.appendLinesWithTabs(1, `if attributes == nil {`)
    writer.appendLinesWithTabs(2, `attributes = map[string]string{}`)
    writer.appendLinesWithTabs(1, `}`)
    writer.appendLinesWithTabs(1, `SetVersionAttribute(${JSON.stringify(message.version)}, attributes)`)
    writer.appendLinesWithTabs(1, `return ${refName}.${engineName}.Publish(${typesPackage}.TopicName${formatEnumName(message.id)}, PubsubMessage{ Data: bytes, Attributes: attributes }, config, ctx)`)
    writer.appendLine( '}')
    writer.appendLine(`func (${refName} *${structName}) RegisterV${message.version}${message.id}Subscription(ff func (data ${typesPackage}.${schema.title}, attributes map[string]string, ctx context.Context) error, config SubscriptionConfig) {`)
    writer.appendLinesWithTabs(1, `${refName}.${schema.title}Registry.RegisterSubscriber(${typesPackage}.TopicName${formatEnumName(message.id)}, config, ff)`)
    writer.appendLine( '}')
}

// The config value is your EventCatalog config
// The options value is the configuration you pass to your plugin
module.exports = async (config, options = {}) => {
    const {
        packageName = 'pubsub',
        preamble = [],
        refName = 'r',
        outPath = 'pubsub_generated.go',
        marshallerName = 'Marshaller',
        engineName = 'Engine',
        structName = 'PubsubRegistry',
        typesPackage = 'types',
    } = options;

    const [commands, events] = await Promise.all([
        getCommands(),
        getEvents(),
        init()
    ])
    const writer = new FileWriter()

    preamble.forEach(p => {
        writer.appendLine(p)
    })

    writer.appendLine(`// Package ${packageName}: Code generated; DO NOT EDIT;`)
    writer.appendLine(`package ${packageName}`)
    writer.appendLine('')
    writer.appendLine(`import (`)
    writer.appendLinesWithTabs(1, `"context"`)
    writer.appendLinesWithTabs(1, `"github.com/staringfun/millsmess/libs/types"`)
    writer.appendLine(')')


    writer.appendLine('')
    writer.appendLine('')

    const commandsToWrite = commands.filter(c => c.channels.some(cc => cc.id !== 'WS'))
    const eventsToWrite = events.filter(e => e.channels.some(cc => cc.id !== 'WS'))

    const messages = [
        ...commandsToWrite,
        ...eventsToWrite,
    ]

    const types = []
    messages.forEach(m => {
        if (!m.schemaPath) return
        const schema = require(`./${m.schemaPath}`)
        types.push(`${schema.title}@${m.channels[0].id}`)
    })
    const allTypes = Array.from(new Set(types))
    writer.appendLine('type PubsubRegistry struct {')
    writer.appendLinesWithTabs(1, 'Marshaller                Marshaller')
    writer.appendLinesWithTabs(1, 'Engine     PubsubEngine')
    allTypes.forEach(t => {
        const tt = t.split('@')[0]
        writer.appendLinesWithTabs(1, `${tt}Registry *TypedSubscribers[${typesPackage}.${tt}]`)
    })
    writer.appendLine('}')

    const channels = {}
    messages.forEach(m => {
        writeMessage(writer, m, channels, {
            refName,
            marshallerName,
            engineName,
            structName,
            typesPackage,
        })
    })
    writer.appendLine(``)
    writer.appendLine(``)
    Object.keys(channels).forEach(c => {
        writer.appendLine(`func (${refName} *${structName}) Handle${c}Message(msg PubsubMessage, config SubscriptionConfig, ctx context.Context) error {`)
        writer.appendLinesWithTabs(1, `version := GetVersionAttribute(msg.Attributes)`)
        writer.appendLinesWithTabs(1, `switch version {`)
        channels[c].events.forEach(e => {
            writer.appendLinesWithTabs(2, `case ${JSON.stringify(e.version)}:`)
            writer.appendLinesWithTabs(3, `var data ${typesPackage}.${e.schemaTitle}`)
            writer.appendLinesWithTabs(3, `err := ${refName}.${marshallerName}.Unmarshal(msg.Data, &data)`)
            writer.appendLinesWithTabs(3, `if err != nil {`)
            writer.appendLinesWithTabs(4, `GetLogger(ctx).Error().Err(err).Msg("unmarshal error")`)
            writer.appendLinesWithTabs(4, `return nil`)
            writer.appendLinesWithTabs(3, `}`)
            writer.appendLinesWithTabs(3, `return ${refName}.${e.schemaTitle}Registry.Run(${typesPackage}.TopicName${formatEnumName(e.id)}, config, data, msg.Attributes, ctx)`)
        })
        writer.appendLinesWithTabs(1, `}`)
        writer.appendLinesWithTabs(1, `if config.IsTopic {`)
        writer.appendLinesWithTabs(2, `return NotMatchedVersionError`)
        writer.appendLinesWithTabs(1, `}`)
        writer.appendLinesWithTabs(1, `return nil`)
        writer.appendLine( '}')
    })
    writer.appendLine(``)
    writer.appendLine(``)
    writer.appendLine(`func (${refName} *${structName}) GetSubscribers() map[types.TopicName]map[SubscriptionConfig]func(PubsubMessage, context.Context) error {`)
    writer.appendLinesWithTabs(1, `topics := map[types.TopicName]map[SubscriptionConfig]func(PubsubMessage, context.Context) error{}`)
    allTypes.forEach(t => {
        const tt = t.split('@')[0]
        const tn = t.split('@')[1]
        writer.appendLinesWithTabs(1, `for _, subscriber := range ${refName}.${tt}Registry.subscribers {`)
        writer.appendLinesWithTabs(2, `if _, ok := topics[subscriber.Topic]; !ok {`)
        writer.appendLinesWithTabs(3, `topics[subscriber.Topic] = map[SubscriptionConfig]func(PubsubMessage, context.Context) error{}`)
        writer.appendLinesWithTabs(2, `}`)
        writer.appendLinesWithTabs(2, `topics[subscriber.Topic][subscriber.Config] = func(msg PubsubMessage, ctx context.Context) error {`)
        writer.appendLinesWithTabs(3, `return ${refName}.Handle${tn}Message(msg, subscriber.Config, ctx)`)
        writer.appendLinesWithTabs(2, `}`)
        writer.appendLinesWithTabs(1, `}`)
    })
    writer.appendLinesWithTabs(1, `return topics`)
    writer.appendLine(`}`)

    const text = writer.getText()
    await fsPromises.writeFile(outPath, format(text))
}