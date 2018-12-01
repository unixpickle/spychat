class ThreadList {
  constructor() {
    this.element = document.getElementById('thread-list');
    this.error = this.element.getElementsByClassName('error')[0];
    this.list = this.element.getElementsByClassName('threads')[0];

    this.onSelectThread = () => null;
    this.abortController = null;
    this.load();
  }

  load() {
    if (this.abortController) {
      this.abortController.abort();
    }
    this.abortController = new AbortController();
    fetch('/threads', { signal: this.abortController.signal })
      .then((response) => response.json())
      .then((jsonResponse) => {
        this.abortController = null;
        if (jsonResponse.error) {
          this.showError(jsonResponse.error);
        } else {
          this.updateThreads(jsonResponse.result);
        }
      })
      .catch(() => {
        this.abortController = null;
        this.showError('request failed')
      });
    this.showLoader();
  }

  showError(message) {
    this.element.className = 'failed';
    this.error.textContent = message;
  }

  showLoader() {
    this.element.className = 'loading';
  }

  updateThreads(threads) {
    this.list.innerHTML = '';

    const threadElems = [];
    threads.forEach((threadInfo, i) => {
      const iconElem = document.createElement('img');
      iconElem.className = 'thread-icon';
      iconElem.src = threadInfoImage(threadInfo);

      const titleElem = document.createElement('label');
      titleElem.className = 'thread-title';
      titleElem.textContent = threadInfoTitle(threadInfo);

      const threadElem = document.createElement('li');
      threadElem.className = 'thread';
      threadElem.appendChild(iconElem);
      threadElem.appendChild(titleElem);
      threadElem.addEventListener('click', () => {
        this.onSelectThread(threadInfo);
        threadElems.forEach((x) => x.className = 'thread');
        threadElem.className = 'thread thread-current';
      });

      threadElems.push(threadElem);
      this.list.appendChild(threadElem);

      if (i === 0) {
        threadElem.className += ' thread-current';
        this.onSelectThread(threadInfo);
      }
    });

    this.element.className = 'loaded';
  }
}

function threadInfoTitle(info) {
  if (info['Name']) {
    return info['Name'];
  }
  let result = 'Group Chat';
  info['Participants'].forEach((p) => {
    if (p['FBID'] == info['OtherUserID']) {
      result = p['Name'];
    }
  });
  return result;
}

function threadInfoImage(info) {
  if (info['Image']) {
    return info['Image'];
  }
  let result = '/assets/svg/no_image.svg';
  info['Participants'].forEach((p) => {
    if (p['FBID'] == info['OtherUserID']) {
      result = p['BigImageSrc'] || p['ImageSrc'];
    }
  });
  return result;
}