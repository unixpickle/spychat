class Messages extends LoadableView {
  constructor() {
    super('messages');
    this.list = this.element.getElementsByClassName('chat-log')[0];
  }

  setThread(thread) {
    this.loadPath('/thread?thread=' + thread['ThreadFBID']);
  }

  showResults(messages) {
    super.showResults(messages);
    this.list.innerHTML = '';
  }
}