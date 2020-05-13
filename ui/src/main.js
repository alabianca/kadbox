const path = require('path');
const kadbox = require('./kadbox')
const { app, BrowserWindow, ipcMain, Tray } = require('electron')


// don't show the app in the dock
//app.dock.hide();
let win;
let tray;

const main = () => {
  createTray();
  createWindow();
};

const createWindow = () => {
  // Create the browser window.
  win = new BrowserWindow({
    width: 320,
    height: 450,
    show: false,
    frame: false,
    fullscreenable: false,
    resizable: false,
    transparent: true,
    webPreferences: {
      nodeIntegration: true
    }
  });

  loadConfig()

  ipcMain.on('window:close', () => win.close());
  ipcMain.on('window:minimize', () => win.minimize());
  ipcMain.on('window:maximize', () => win.isMaximized() ? win.unmaximize() : win.maximize());

  win.webContents.openDevTools()

  // and load the index.html of the app.
  win.loadURL('http://localhost:3000');

  // Hide the window when it loses focus
  win.on('blur', () => {
    if (!win.webContents.isDevToolsOpened()) {
      win.hide();
    }
  });
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


const loadConfig = async () => {
  try {
    const config = await kadbox.loadConfig()
    //console.log(config)
    kadbox.box();
    await kadbox.ping();
  } catch (e) {
    console.log("Error ", e)
  }
}


app.whenReady().then(main)