class Messages extends LoadableView {
  constructor() {
    super('messages');
    this.list = this.element.getElementsByClassName('chat-log')[0];
    this.thread = null;
  }

  setThread(thread) {
    this.thread = thread;
    this.loadPath('/thread?thread=' + thread['ThreadFBID']);
  }

  showResult(messages) {
    super.showResult(messages);

    this.list.innerHTML = '';
    messages.forEach((msg) => {
      if (!msg['Body'] && !msg['Attachments']) {
        return;
      }

      const info = this.userInfo((msg['RawData']['message_sender'] || {})['id']);

      const iconElem = document.createElement('img');
      iconElem.className = 'message-icon';
      iconElem.src = info.icon;

      const nameElem = document.createElement('label');
      nameElem.className = 'message-name';
      nameElem.textContent = info.name;

      const bodyElem = document.createElement('div');
      bodyElem.className = 'message-body';
      bodyElem.textContent = msg['Body'];

      const messageElem = document.createElement('div');
      messageElem.className = 'message';
      messageElem.appendChild(iconElem);
      messageElem.appendChild(nameElem);
      messageElem.appendChild(bodyElem);

      this.list.appendChild(messageElem);
    });

    this.element.scrollTo(0, this.list.offsetHeight);
  }

  userInfo(fbid) {
    let result = { name: 'Unknown', icon: '/assets/svg/no_image.svg' };
    this.thread['Participants'].forEach((p) => {
      if (p['FBID'] == fbid) {
        result.name = p['Name'];
        result.icon = p['BigImageSrc'] || p['ImageSrc'];
      }
    });
    return result;
  }
}