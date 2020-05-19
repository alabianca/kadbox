import { useState, useEffect } from 'react';
import Lib from '../lib'
import {IpcEventBus} from "../lib/eventBus";

export const useNewFileCounter = (initial: number = 0) => {
    const [counter, setCounter] = useState(initial);
    const cmgr = Lib.getConnectionManager();

    useEffect(() => {
        const sub = cmgr.EventBus().subscribe(IpcEventBus.STORAGE_SUCCESS, () => {
            console.log("Got a new file")
            setCounter(counter + 1);
        })

        return function cleanup() {
            sub.unsubscribe();
        }
    })

    return {
        counter,
        setCounter,
    }
}