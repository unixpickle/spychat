class App {
  constructor() {
    this.threadList = new ThreadList();
    this.messages = new Messages();
  }
}

window.addEventListener('load', () => window.app = new App());