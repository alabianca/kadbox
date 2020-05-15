const homedir = require('os').homedir();
const { promises } = require('fs');
const path = require('path');
const { spawn } = require('child_process');
const nodeFetch = require('node-fetch')

const STORE_LOC = "store";
const CONFIG_LOC = "kadconfig";
const ROOT_LOC = ".kadbox";

const loadConfig = async () => {
    const data = await promises.readFile(path.join(homedir, ROOT_LOC, CONFIG_LOC));
    //console.log(data.toString())

    return JSON.parse(data.toString())
}

const box = () => {
    const childproc = spawn('kadbox', ['server', 'start']);
    childproc.stdout.pipe(process.stdout);

    return childproc;
}

const ping = async () => {
    const config = await loadConfig();
    const url = `http://${config.api.address}:${config.api.port}/ping`
    await nodeFetch(url)
}


module.exports = {
    STORE_LOC,
    ROOT_LOC,
    CONFIG_LOC,
    loadConfig,
    box,
    ping,
}