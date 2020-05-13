
const getIpcRenderer = () => {
    if (!window || !window.process || !window.require) {
        throw new Error("IPC Renderer is not available")
    }

    return window.require('electron').ipcRenderer;
}

export default getIpcRenderer;