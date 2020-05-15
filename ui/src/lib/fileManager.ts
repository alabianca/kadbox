import {Dispatcher} from "./dispatcher";
import {EVENT_ADD_FILE} from './connectionManager'

export class FileManager {
    constructor(private dispatcher: Dispatcher) {
    }

    public addObject(objectPath: string) {
        this.dispatcher.dispatch(EVENT_ADD_FILE, objectPath)
    }
}