const { ipcMain } = require('electron');
const { EVENT_PROCESS_PING, EVENT_PROCESS_PONG, EVENT_PROCESS_KILL } = require('./ipcManager')

class ProcessManager {
    constructor(proc) {
        console.log("Process Manager is running")
        this._proc = proc;
        this.running = true;
        // ipcMain.on('kad:ping', (event) => {
        //     if (this.running) {
        //         event.sender.send("kad:pong")
        //     } else {
        //         event.sender.send("kad:kill")
        //     }
        // });

        this._proc.on('exit', () => {
            console.log('Kadbox exited')
            this.running = false;
        })
    }

    dispatch(name, event, args) {
        switch (name) {
            case EVENT_PROCESS_PING:
                this._pong(event);
                break;

            default:
                console.log("Action not supported by ProcessManager")
        }
    }

    _pong(event) {
        if (this.running) {
            event.sender.send(EVENT_PROCESS_PONG)
        } else {
            event.sender.send(EVENT_PROCESS_KILL);
        }
    }

    setProc(proc) {
        this._proc = proc;
        this.running = true;
    }
}

module.exports = {
    ProcessManager,
}