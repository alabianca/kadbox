const { EVENT_STORAGE_ADD, EVENT_STORAGE_SUCCESS, EVENT_STORAGE_ERROR, EVENT_STORAGE_GET_FILES, EVENT_STORAGE_FILES } = require('./ipcManager')
const request = require('request');
const tar = require('tar-fs');
const { createGzip } = require('zlib');
const { pipeline, Duplex } = require('stream');
const { createWriteStream, createReadStream, readdirSync } = require('fs')
const kadbox = require('./kadbox')

class StorageManager {
    constructor() {
    }

    dispatch(name, event, arg) {
        switch (name) {
            case EVENT_STORAGE_ADD:
                this._addObject(event, arg);
                break;
            case EVENT_STORAGE_GET_FILES:
                this._getObjects(event);
        }
    }

    async _addObject(event, objectPath) {
        try {
            const config = await kadbox.loadConfig();
            const url = `http://${config.api.address}:${config.api.port}/storage`;
            await uploadFile(objectPath, url);
            event.sender.send(EVENT_STORAGE_SUCCESS);
        } catch (e) {
            event.sender.send(EVENT_STORAGE_ERROR)
        }
    }

     _getObjects(event) {
        const packedFiles = readdirSync(kadbox.storeDir());
        event.sender.send(EVENT_STORAGE_FILES, packedFiles)
    }
}

const uploadFile = async (filePath, url) => {
    try {
        const tempFile = "./tmp.tar.gz";
        await packAndCompress(filePath, tempFile);
        const formData = {
            upload: createReadStream(tempFile),
        }
        request.post({
            url,
            formData,
        }, (err, res, body) => {
            if (err) {
                console.log("Error", err)
                throw new Error(err)

            }
            console.log(body)
        })
    } catch (e) {
        console.log("Upload Error", e)
    }
}

const packAndCompress = async (filePath, tempFileName) => {

    const source = tar.pack(filePath);
    return new Promise((resolve, reject) => {
        pipeline(source, createGzip(), createWriteStream(tempFileName), (err) => {
            if (err) {
                reject(err)
            }
            resolve()
        })
    })
}

module.exports = {
    StorageManager,
}