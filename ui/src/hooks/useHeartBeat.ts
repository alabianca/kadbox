import {useEffect, useState} from "react";
import Lib from "../lib"
import {IpcEventBus} from "../lib/eventBus";

export const useHeartBeat = () => {
    const cmgr = Lib.getConnectionManager();
    const [isConnected, setIsConnected] = useState(false);

    useEffect(() => {
        const subPong = cmgr.EventBus().subscribe(IpcEventBus.PONG, () => setIsConnected(true));
        const subKill = cmgr.EventBus().subscribe(IpcEventBus.KILLED, () => setIsConnected(false))
        // ping every 5 seconds to make sure we are still connected to the network
        const ivalId = setInterval(() => {
            cmgr.ping();
        }, 10000)

        return function cleanup() {
            subPong.unsubscribe();
            subKill.unsubscribe();
            clearInterval(ivalId);
        }
    });

    return {
        isConnected,
    }
}