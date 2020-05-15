const { EVENT_STORAGE_ADD } = require('./ipcManager')

class StorageManager {
    constructor() {
    }

    dispatch(name, event, arg) {
        switch (name) {
            case EVENT_STORAGE_ADD:
                this._addObject(arg)
        }
    }

    async _addObject(objectPath) {
        console.log('Adding ', objectPath)
    }
}

module.exports = {
    StorageManager,
}