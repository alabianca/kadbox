const { ipcMain } = require('electron');

const EVENT_WINDOW_CLOSE = "window:close";
const EVENT_WINDOW_MAXIMIZE = "window:maximize";
const EVENT_WINDOW_MINIMIZE = "window:minimize";
const EVENT_PROCESS_PING = "kad:ping";
const EVENT_PROCESS_PONG = "kad:pong";
const EVENT_PROCESS_KILL = "kad:kill";
const EVENT_STORAGE_ADD = "storage:add";
const EVENT_STORAGE_SUCCESS = "storage:success";
const EVENT_STORAGE_ERROR = "storage:error";
const EVENT_STORAGE_FILES = "storage:files";
const EVENT_STORAGE_GET_FILES = "storage:get-files"

class IpcManager {
    constructor() {
        this._dispatchers = {};

        // window events
        ipcMain.on(EVENT_WINDOW_CLOSE, (event, arg) => this._locateAndDispatch(EVENT_WINDOW_CLOSE, event, arg));
        ipcMain.on(EVENT_WINDOW_MINIMIZE, (event, arg) => this._locateAndDispatch(EVENT_WINDOW_MINIMIZE, event, arg));
        ipcMain.on(EVENT_WINDOW_MAXIMIZE, (event, arg) => this._locateAndDispatch(EVENT_WINDOW_MAXIMIZE, event, arguments));
        // process manager events
        ipcMain.on(EVENT_PROCESS_PING, (event, arg) => this._locateAndDispatch(EVENT_PROCESS_PING, event, arg));
        ipcMain.on(EVENT_PROCESS_KILL, (event, arg) => this._locateAndDispatch(EVENT_PROCESS_KILL, event, arg));
        // storage events
        ipcMain.on(EVENT_STORAGE_ADD, (event, arg) => this._locateAndDispatch(EVENT_STORAGE_ADD, event, arg));
        ipcMain.on(EVENT_STORAGE_GET_FILES, (event, arg) => this._locateAndDispatch(EVENT_STORAGE_GET_FILES, event, arg))
    }

    addDispatcher(key, dispatcher) {
        if (this._dispatchers[key]) {
            this._dispatchers[key].push(dispatcher)
        } else {
            this._dispatchers[key] = [dispatcher];
        }
    }

    _locateAndDispatch(name, event, arg) {
        const dispKey = name.split(':')[0];
        const dispatchers = this._dispatchers[dispKey];
        if (dispatchers) {
            dispatchers.forEach((d) => d.dispatch(name, event, arg))
        }
    }
}

module.exports = {
    IpcManager,
    EVENT_PROCESS_KILL,
    EVENT_PROCESS_PING,
    EVENT_PROCESS_PONG,
    EVENT_STORAGE_ADD,
    EVENT_STORAGE_ERROR,
    EVENT_STORAGE_SUCCESS,
    EVENT_STORAGE_GET_FILES,
    EVENT_STORAGE_FILES,
}