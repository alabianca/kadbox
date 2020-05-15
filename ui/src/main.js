const path = require('path');
const kadbox = require('./etron/kadbox')
const { app, BrowserWindow, ipcMain, Tray } = require('electron');
const { sleep } = require('./etron/utils')
const { ProcessManager } = require('./etron/processManager');
const { IpcManager } = require('./etron/ipcManager');
const { StorageManager } = require('./etron/storageManager');

// don't show the app in the dock
//app.dock.hide();
let win;
let tray;

const main = () => {
  createTray();
  createWindow();
};

const createWindow = async () => {
  // Create the browser window.
  win = new BrowserWindow({
    width: 380,
    height: 550,
    show: false,
    frame: false,
    fullscreenable: false,
    resizable: false,
    transparent: true,
    webPreferences: {
      nodeIntegration: true
    }
  });

  win.webContents.openDevTools()

  // and load the index.html of the app.
  win.loadURL('http://localhost:3000');

  // Hide the window when it loses focus
  win.on('blur', () => {
    if (!win.webContents.isDevToolsOpened()) {
      win.hide();
    }
  });

  const ipcManager = new IpcManager();
  const storage = new StorageManager();
  ipcManager.addDispatcher('storage', storage);

  try {
    const proc = await startKadboxProcess();
    const pm = new ProcessManager(proc);
    ipcManager.addDispatcher('kad', pm);

  } catch (e) {
  }
};

const createTray = () => {
  tray = new Tray(path.resolve(__dirname, '../public', 'eth.png'));
  tray.on('click', toggleWindow)
};

const toggleWindow = () => {
  win.isVisible() ? win.hide() : showWindow();
};

const showWindow = () => {
  const pos = getWindowPosition();
  win.setPosition(pos.x, pos.y, false);
  win.show();
};

const getWindowPosition = () => {
  const windowBounds = win.getBounds();
  const trayBounds = tray.getBounds();

  // Center window horizontally below the tray icon
  const x = Math.round(trayBounds.x + (trayBounds.width / 2) - (windowBounds.width / 2));
  // Position window 4 pixels vertically below the tray icon
  const y = Math.round(trayBounds.y + trayBounds.height + 4);
  return {x: x, y: y}
}


const startKadboxProcess = async (attempt = 0, proc = null) => {
  if (attempt > 3) {
    throw new Error("cannot start kadbox process")
  }
  try {
    // first check if it is up already
    await kadbox.ping();
    console.log(`kadbox is running`)
  } catch (e) {
    console.log(`kadbox start up attempt ${attempt + 1} out of 3`);
    proc = kadbox.box();
    attempt++;
    await sleep(2000); // give it a 2 seconds to start up before trying again
    await startKadboxProcess(attempt, proc);
  }

  return proc;
}


app.whenReady().then(main)