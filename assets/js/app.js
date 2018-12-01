class App {
  constructor() {
    this.threadList = new ThreadList();
    this.messages = new Messages();
    this.threadList.onSelectThread = (th) => this.messages.setThread(th);
  }
}

window.addEventListener('load', () => window.app = new App());