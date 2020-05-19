import {Dispatcher} from "./dispatcher";
import {EVENT_ADD_FILE, EVENT_GET_FILE} from './connectionManager'
import {IpcEventBus, Subscriber} from './eventBus'

export class FileManager {
    constructor(private dispatcher: Dispatcher, private subscriber: Subscriber) {
    }

    public addObject(objectPath: string) {
        this.dispatcher.dispatch(EVENT_ADD_FILE, objectPath)
    }

    public getFiles(callback: (files: string[]) => any) {
        const sub = this.subscriber.subscribe(IpcEventBus.STORAGE_FILES, (files: string[]) => {
            callback(files);
            sub.unsubscribe();
        })

        this.dispatcher.dispatch(EVENT_GET_FILE)
    }
}