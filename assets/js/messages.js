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
      bodyElem.innerHTML = escapeBody(msg['Body']);

      (msg['Attachments'] || []).forEach((attachment) => {
        let attachElem = null;
        if (attachment['HiResURL']) {
          attachElem = document.createElement('img');
          attachElem.className = 'message-attachment-image';
          attachElem.src = attachment['HiResURL'];
        } else if (attachment['AudioURL']) {
          attachElem = document.createElement('audio');
          attachElem.className = 'message-attachment-audio';
          attachElem.controls = 'controls';
          attachElem.src = attachment['AudioURL'];
        } else if (attachment['VideoURL']) {
          attachElem = document.createElement('video');
          attachElem.className = 'message-attachment-video';
          attachElem.controls = 'controls';
          attachElem.src = attachment['VideoURL'];
        } else if (attachment['FileURL']) {
          attachElem = document.createElement('a');
          attachElem.className = 'message-attachment-file';
          attachElem.href = attachment['FileURL'];
          attachElem.textContent = 'Download file (' + attachment['Name'] + ')';
          attachElem.target = '_blank';
        } else if (attachment['StickerID']) {
          attachElem = document.createElement('img');
          attachElem.className = 'message-attachment-sticker';
          attachElem.src = attachment['RawURL'];
        }
        if (attachElem) {
          bodyElem.appendChild(attachElem);
          attachElem.onload = () => {
            this.element.scrollTo(0, this.list.offsetHeight);
          }
        }
      });

      const messageElem = document.createElement('div');
      messageElem.className = 'message';
      messageElem.appendChild(iconElem);
      messageElem.appendChild(nameElem);
      messageElem.appendChild(bodyElem);

      this.list.appendChild(messageElem);
    });

    this.element.scrollTo(0, this.list.offsetHeight + 100000);
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

function escapeBody(body) {
  body = body.replace(/&/g, '&amp;');
  body = body.replace(/</g, '&lt;');
  body = body.replace(/>/g, '&gt;');
  body = body.replace(/\n/g, '<br>');
  return body;
}