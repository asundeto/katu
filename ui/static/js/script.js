//WRAPPER MAIN ELEMENTS
// const box = document.querySelector('.box')
// const create = document.querySelector('.create')
// const post = document.querySelector('.post_main')
const wrapper = document.querySelector('.wrapper')
const loginLink = document.querySelector('.login-link')
const registerLink = document.querySelector('.register-link')

//ADAPTIVE MENU ELEMENTS
const menu = document.querySelector(".icon-menu")
const navLink = document.querySelector(".navigation")
const btnPop = document.querySelector('.btnLogin')
const navBack = document.getElementById('nav-back-mobile')

btnPop.addEventListener('click', () => {
  if (navLink.classList.contains('mob-menu')) {
    navLink.classList.remove('mob-menu')
  }
  if (wrapper.classList.contains('active-pop')) {
    wrapper.classList.remove('active-pop')
  } else {
    wrapper.classList.add('active-pop')
  }
});

menu.addEventListener('click', () => {
  navLink.classList.toggle('mob-menu')
  if (navBack.style.display == "block") {
    navBack.style.display = "none"
  } else {
    navBack.style.display = "block"
  }
})

navBack.addEventListener('click', () => {
  navLink.classList.toggle('mob-menu')
  navBack.style.display = "none"
})

if (registerLink != null && loginLink != null) {
  registerLink.addEventListener('click', () => {
    wrapper.classList.add('active')
    wrapper.classList.remove('error')
  });
  
  loginLink.addEventListener('click', () => {
    wrapper.classList.remove('active')
    wrapper.classList.remove('error')
  });
}

// -------------------- Check Login Parametrs ------------------- //
//PASSWORD LOCK ICONS 
// const logIconLock1 = document.querySelector('.log-icon-lock1')
// const logIconLock2 = document.querySelector('.log-icon-lock2')
// const logIconLock3 = document.querySelector('.log-icon-lock3')

// const locks = [0, 0, 0];

// function toggleLock(lockIndex, iconElement, passwordElement) {
//   const lock = locks[lockIndex];
//   if (lock % 2 === 0) {
//     iconElement.setAttribute('name', 'lock-open');
//     passwordElement.setAttribute('type', 'text');
//   } else {
//     iconElement.setAttribute('name', 'lock-closed');
//     passwordElement.setAttribute('type', 'password');
//   }
//   locks[lockIndex]++;
// }

//PASSWORD LABELS
const password_lbl_log = document.getElementById('password-label-log')
const password_lbl_reg = document.getElementById('password-label-reg')
const password_lbl_reg2 = document.getElementById('password-label-reg2')

// if (logIconLock1 && logIconLock2 && logIconLock3) {
//   logIconLock1.addEventListener('click', () => toggleLock(0, logIconLock1, password_lbl_log));
//   logIconLock2.addEventListener('click', () => toggleLock(1, logIconLock2, password_lbl_reg));
//   logIconLock3.addEventListener('click', () => toggleLock(2, logIconLock3, password_lbl_reg2));
// }

// -------------------- Check Login Parametrs ------------------- //

//INPUTS
const logEmail = document.querySelector('.log-email')
const logPassword = document.querySelector('.log-pass')
const regUsername = document.querySelector('.reg-username')
const regEmail = document.querySelector('.reg-email')
const regPassword = document.querySelector('.reg-password')
const regPasswordSec = document.querySelector('.reg-password-second')
const checkbox = document.getElementById('register-checkbox')

const error_window_box = document.getElementById("error-window-box")
const error_window_box2 = document.getElementById("error-window-box2")
const error_window_log = document.getElementById("error-window-log")
const error_window_log2 = document.getElementById("error-window-log2")

//GET FROM BACK INPUT VALUES
const ReturnedEmail = document.getElementById('ReturnedEmail')
const ReturnedUsername = document.getElementById('ReturnedUsername')
const ReturnedPass = document.getElementById('ReturnedPass')
const ReturnedPass2 = document.getElementById('ReturnedPass2')

//TO GET ERRORS FROM BACK
const form_act = document.getElementById("form-act")
const form_error = document.getElementById("form-error")

// -------------------- Image upload Drop Down Area ------------------- //
const dropArea = document.getElementById('drop-area');
const fileInput = document.getElementById('fileInput');

if (dropArea) {
  dropArea.addEventListener('dragover', (e) => {
    e.preventDefault();
    dropArea.style.border = '2px dashed #28a745';
  });

  dropArea.addEventListener('dragleave', () => {
    dropArea.style.border = '2px dashed #ccc';
  });

  dropArea.addEventListener('drop', (e) => {
    e.preventDefault();
    dropArea.style.border = '2px dashed #ccc';

    const files = e.dataTransfer.files;

    if (files.length > 0) {
      fileInput.files = files;
      handleFiles(files);
    }
  });

  fileInput.addEventListener('change', () => {
    const files = fileInput.files;
    handleFiles(files);
  });

  function handleFiles(files) {
    for (const file of files) {
      if (file.type.startsWith('image/')) {
        const imgElement = document.getElementById("show-uploaded-image");
        imgElement.src = URL.createObjectURL(file);
        imgMessage.innerHTML = "Image upload success";
        imgMessage.style.display = "block";
      } else {
        imgMessage.innerHTML = "Invalid file type!";
        imgMessage.style.display = "block";
      }
    }
  }
}

// ------------------------- MESSAGES -------------------------
const imgMessage = document.querySelector('.post-image-message');
const allMessages = document.querySelectorAll('.message-text');
const lastMessage = document.querySelectorAll('.chat-last-message');
const userRow = document.getElementById('user-row');
const myChatsBox = document.getElementById('my-chats-box')
const chatMessageBox = document.getElementById('chat-message-box')

if (lastMessage) {
  shortMessage()
}

function shortMessage() {
  for (let i = 0; i < lastMessage.length; i++) {
    lastMessage[i].innerHTML = shortLastMessage(removeSpaces(lastMessage[i].innerHTML))
  }
}

function removeSpaces(str) {
  var result = ""
  for (let i = 0; i < str.length; i++) {
    if (str[i] != ' ') {
      result += str[i]
    }
  }
  return result
}

function shortLastMessage(message) {
  var shorted = ""
  var width = myChatsBox.offsetWidth
  var max = Math.floor(width / 21)
  var cnt = 0
  var breaked = false
  for (let i = 0; i < message.length; i++) {
    if (cnt == max) {
      breaked = true
      break
    }
    shorted += message[i]
    cnt++
  }
  if (breaked) {
    return shorted + "..."
  } else {
    return shorted
  }
}

const createTypesBtn = document.getElementById('create-types-btn')
const createTypesContainer = document.getElementById('create-types-container') 
if (createTypesBtn != null) {
  let clickCount = 0;
  createTypesBtn.addEventListener('click', () => {
    clickCount++;
    if (clickCount % 2 !== 0) {
        createTypesContainer.style.display = "flex";
    } else {
        createTypesContainer.style.display = "none";
    }
});
}


// ???????????????????????????????????????????????????????????????????

// -------------------- CHAT ------------------- //
const inputElementChat = document.getElementById("chat-input");
if (inputElementChat != null) {
  window.onload = function() {
    scrollToBottom();
    inputElementChat.focus();
  };
}

// -------------------- Chat Auto Scroll Down ------------------- //
let msgChatScrollHeight = 0
const messagesChatWindow = document.getElementById('messagesChatWindow')

function scrollToBottom() {
  if (messagesChatWindow) {
    messagesChatWindow.scrollTop = messagesChatWindow.scrollHeight;
  }
  msgChatScrollHeight = messagesChatWindow.scrollTop
}

function scrollDown() {
  let msgChatScrollHeight2 = 0
  document.getElementById('messagesChatWindow').scrollTop = document.getElementById('messagesChatWindow').scrollHeight;
  msgChatScrollHeight2 = document.getElementById('messagesChatWindow').scrollTop
}

// -------------------- Chat Enter Press Event ------------------- //

function handleKeyDown(event) {
  if (event.keyCode === 13) {
    sendMessage();
    document.getElementById("chat-input").focus();
  }
}

// -------------------- Change value of chat window ------------------- //
function extractChatWindow(htmlString) {
  let tempElement = document.createElement('div');
  tempElement.innerHTML = htmlString;
  
  let chatWindowDiv = tempElement.querySelector('div.messages-chat-window#messagesChatWindow');

  if (chatWindowDiv) {
      return chatWindowDiv.innerHTML;
  } else {
      return null;
  }
}

// -------------------- Get Category Func ------------------- //
function getCategoryCode(data) {
  let tempElement = document.createElement('div');
  tempElement.innerHTML = data;
  
  let categiryWindowDiv = tempElement.querySelector('div.post_table#post-table');

  if (categiryWindowDiv) {
      return categiryWindowDiv.innerHTML;
  } else {
      return null;
  }
}

function getPostData(data) {
  let tempElement = document.createElement('div');
  tempElement.innerHTML = data;
  
  let postReaction = tempElement.querySelector('div.post-box-reaction-back');

  if (postReaction) {
      return postReaction.innerHTML;
  } else {
      return null;
  }
}

function getCommentData(data) {
  let tempElement = document.createElement('div');
  tempElement.innerHTML = data;
  
  let commentReeaction = tempElement.querySelector('div.comments');

  if (commentReeaction) {
      return commentReeaction.innerHTML;
  } else {
      return null;
  }
}

// -------------------- Register Func ------------------- //
function getRegisterRequest(data) {
  let tempElement = document.createElement('div');
  tempElement.innerHTML = data;

  let categiryWindowDiv = tempElement.querySelector('div.post_table#post-table');

  if (categiryWindowDiv) {
      return categiryWindowDiv.innerHTML;
  } else {
      return null;
  }
}

function dataConverter(data, clas, id) {
  let tempElement = document.createElement('div');
  tempElement.innerHTML = data;
  let result

  if (id == null) {
    result = tempElement.querySelector('div.'+clas);
  } else {
    result = tempElement.querySelector('div.'+clas+"#"+id);
  }

  if (result) {
      return result.innerHTML;
  } else {
      return null;
  }
}

function dataConverterChat(data, clas, id) {
  let tempElement = document.createElement('div');
  tempElement.innerHTML = data;
  let result

  if (id == null) {
    result = tempElement.querySelector('div.'+clas);
  } else {
    result = tempElement.querySelector('div.'+clas+"#"+id);
  }

  if (result) {
      return result.innerHTML;
  } else {
      return null;
  }
}

function dataConverterLastMessage(data, clas, id) {
  let tempElement = document.createElement('div');
  tempElement.innerHTML = data;
  let result

  if (id == null) {
    result = tempElement.querySelector('div.'+clas);
  } else {
    result = tempElement.querySelector('div.'+clas+"#"+id);
  }
  var allChats = result.querySelectorAll('.chat-last-message')
  for (let i = 0; i < allChats.length; i++) {
    allChats[i].innerHTML = shortLastMessage(removeSpaces(allChats[i].innerHTML))
  }

  if (result) {
      return result.innerHTML;
  } else {
      return null;
  }
}

// const msgChatWindow = document.getElementById('messagesChatWindow')
const messagesPage = document.getElementById('messages-box-back')
if (messagesPage != null) {
  //refreshChat()
} 


const scrollDownBtn = document.getElementById('scrollDownBtn')

function debounce(func, delay) {
  let timeoutId;
  return function() {
    const context = this;
    const args = arguments;
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => {
      func.apply(context, args);
    }, delay);
  };
}

//Debounced scroll event listener

function scrollDownFunction() {
  let lastScrollTop = 0;
  document.getElementById('messagesChatWindow').addEventListener("scroll", function() {
    let scrollPosition = document.getElementById('messagesChatWindow').scrollTop;
      
    if (msgChatScrollHeight - scrollPosition > 20) {
      document.getElementById('scrollDownBtn').style.display = "block";
      document.getElementById('scrollDownBtn').style.animation = 'dropDownBtnIn 0.5s forwards';
    } else {
      document.getElementById('scrollDownBtn').style.animation = 'dropDownBtnOut 0.5s forwards';
      setTimeout(() => {
        document.getElementById('scrollDownBtn').style.display = "none";
      }, 500);
    }
    document.getElementById('scrollDownBtn').addEventListener('click', () => {
      scrollDown()
    })
  });
}
if (document.getElementById('messagesChatWindow')) {
  scrollDownFunction()
}


// if (messagesChatWindow) {
//   messagesChatWindow.addEventListener("scroll", debounce(function() {
//     let scrollPosition = messagesChatWindow.scrollTop;
      
//     if (msgChatScrollHeight - scrollPosition > 20) {
//       scrollDownBtn.style.display = "block";
//       scrollDownBtn.style.animation = 'dropDownBtnIn 0.5s forwards';
//     } else {
//       scrollDownBtn.style.animation = 'dropDownBtnOut 0.5s forwards';
//       setTimeout(() => {
//         scrollDownBtn.style.display = "none";
//       }, 500);
//     }
//   }, 200));
// }

if (scrollDownBtn != null) {
  scrollDownBtn.addEventListener('click', () => {
    scrollToBottom()
  })
}

// -------------------- Login -----------------------------

function login() {
  var badRequest = false
  fetch("/user/login/post", {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded"
    },
    // body: formData,
    body: "email=" + encodeURIComponent(logEmail.value) + "&password=" + encodeURIComponent(logPassword.value),
  })
    .then(response => {
      if (!response.ok) {
        if (response.status == 400) {
          badRequest = true
        } else {
          throw new Error("Network response was not ok");
        }
    }
    return response.text();
    })
    .then((data) => {
      if (badRequest) {
        const parsedJson = JSON.parse(data);
        document.getElementById('login-error').innerHTML = parsedJson.Error
        return
      }
      document.documentElement.innerHTML = data
      history.pushState(null, "", "/");
    })
    .catch((error) => {
      console.error("Error:", error);
    });
}

// -------------------- Register -----------------------------

function register() {
  var badRequest = false
  if (checkbox.checked) {
    fetch("/user/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      },
      // body: formData,
      body: "username=" + encodeURIComponent(regUsername.value) + "&email=" + encodeURIComponent(regEmail.value) + "&password=" + encodeURIComponent(regPassword.value) + "&password-repeat=" + encodeURIComponent(regPasswordSec.value) + "&checkbox=" + encodeURIComponent(checkbox.value),
    })
      .then(response => {
        if (!response.ok) {
          if (response.status == 400) {
            badRequest = true
          } else {
            throw new Error("Network response was not ok");
          }
      }
      return response.text();
      })
      .then((data) => {
        if (badRequest) {
          const parsedJson = JSON.parse(data);
          document.getElementById('register-error').innerHTML = parsedJson.Error
          return
        }
        regUsername.value = ""
        regEmail.value = ""
        regPassword.value = ""
        regPasswordSec.value = ""
        checkbox.checked = false
        wrapper.classList.remove('active')
        document.getElementById('register-error').innerHTML = ""
      })
      .catch((error) => {
        console.error("Error:", error);
      }); 
  } else {
    document.getElementById('register-error').innerHTML = "Accept user agreement!"
  }
}

// -------------------- Create Post -----------------------

const titleError = document.getElementById('title-error')
const contentError = document.getElementById('content-error')
const imageError = document.getElementById('image-error')

function createPost() {
  const form = document.getElementById("createPostForm");
  const formData = new FormData(form);
  var badRequest = false
  fetch("/post/create/post", {
    method: "POST",
    body: formData,
  })
    .then(response => {
      if (!response.ok) {
        if (response.status == 400) {
          badRequest = true
        } else {
          throw new Error("Network response was not ok");
        }
    }
    return response.text();
    })
    .then((data) => {
      if (badRequest) {
        const parsedJson = JSON.parse(data);
        titleError.innerHTML = parsedJson.TitleError
        contentError.innerHTML = parsedJson.ContentError
        imageError.innerHTML = parsedJson.ImageError
        return
      }
      document.documentElement.innerHTML = data
    })
    .catch((error) => {
      console.error("Error:", error);
    });
}

const messagesUserName = document.getElementById('user-messages-username')

// ----------------------- Open Chat ---------------------

async function openChat(withName) {
  if (screen.width < 1170) {
    document.getElementById('my-chats-box').style.display = 'none'
    document.getElementById('chat-message-box').style.display = 'flex'
  }
  history.pushState(null, "", "/chat/" + withName);
    fetch("/chat/"+withName, {
        method: "GET",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded"
        },
    })
    .then(response => {
        if (!response.ok) {
            throw new Error("Network response was not ok");
        }
        return response.text();
    })
    .then(data => {
      if(chatMessageBox) {
        var chat_updated = dataConverter(data, "chat-message-box", "chat-message-box")
        if (chatMessageBox.innerHTML != chat_updated) {
          chatMessageBox.innerHTML = chat_updated
          let msgChatScrollHeight = 0
          const messagesChatWindow1 = document.getElementById('messagesChatWindow')
          messagesChatWindow1.scrollTop = messagesChatWindow1.scrollHeight;
          msgChatScrollHeight = messagesChatWindow1.scrollTop
          document.getElementById("chat-input").focus();
          funcToChatImageUpload()
          chatImageDropFunction()
          activateChatBackBtn()
          scrollDown()
        }
      }
    })
    .catch(error => {
        console.error("There was a problem with the fetch operation:", error);
    });
}

function refreshChat() {
  setInterval(() => {
    var withName = messagesUserName != null ? messagesUserName.innerHTML : "";
    fetch("/chat/"+withName, {
        method: "GET",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded"
        },
    })
    .then(response => {
        if (!response.ok) {
            throw new Error("Network response was not ok");
        }
        return response.text();
    })
    .then(data => {
      if(messagesChatWindow) {
        var chat_updated = dataConverter(data, "messages-chat-window", "messagesChatWindow")
        if (messagesChatWindow.innerHTML != chat_updated) {
          messagesChatWindow.innerHTML = chat_updated
          // scrollToBottom();
        }
      }
      var my_chats_updated = dataConverterLastMessage(data, "my-chats-box", "my-chats-box")
      if (document.getElementById("my-chats-box").innerHTML != my_chats_updated) {
        document.getElementById("my-chats-box").innerHTML = my_chats_updated
      }
      scrollDown()
      //shortMessage()
    })
    .catch(error => {
        console.error("There was a problem with the fetch operation:", error);
    });
  }, 1000)
}

// -------------------- RefreshChatOnce ---------------------------//

function refreshChatOnce() {
  var withName = messagesUserName != null ? messagesUserName.innerHTML : "";
  history.pushState(null, "", "/chat/" + withName);
    fetch("/chat/"+withName, {
        method: "GET",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded"
        },
    })
    .then(response => {
        if (!response.ok) {
            throw new Error("Network response was not ok");
        }
        return response.text();
    })
    .then(data => {
      if(messagesChatWindow) {
        var chat_updated = dataConverter(data, "messages-chat-window", "messagesChatWindow")
        if (messagesChatWindow.innerHTML != chat_updated) {
          messagesChatWindow.innerHTML = chat_updated
        }
      }
      var my_chats_updated = dataConverterLastMessage(data, "my-chats-box", "my-chats-box")
      if (document.getElementById("my-chats-box").innerHTML != my_chats_updated) {
        document.getElementById("my-chats-box").innerHTML = my_chats_updated
      }
      runSwitcher()
      scrollDown()
    })
    .catch(error => {
        console.error("There was a problem with the fetch operation:", error);
    });
}

// -------------------- Chat Send Message Func ------------------- //

// const chatUploadImage = document.getElementById('chatUploadImage');
// const chatFileUploader = document.getElementById('chatFileUploader');
// const chatImageShowBack = document.getElementById('chatUplodedImageShowBack')


// const chatDropArea = document.getElementById('messagesChatWindow');
// const chatFileInput = document.getElementById('fileInput');

if (document.getElementById('messagesChatWindow')) {
  chatImageDropFunction()
}

if (document.getElementById('chatUploadImage')) {
  funcToChatImageUpload()
}

function chatImageDropFunction() {
  document.getElementById('messagesChatWindow').addEventListener('dragover', (e) => {
    e.preventDefault();
    document.getElementById('messagesChatWindow').style.border = '2px dashed #2C1B7D';
  });

  document.getElementById('messagesChatWindow').addEventListener('dragleave', () => {
    document.getElementById('messagesChatWindow').style.border = '2px solid #30363d';
  });

  document.getElementById('messagesChatWindow').addEventListener('drop', (e) => {
    e.preventDefault();
    document.getElementById('messagesChatWindow').style.border = '2px solid #30363d';

    const files = e.dataTransfer.files;

    if (files.length > 0) {
      document.getElementById('chatFileUploader').files = files;
      handleChatUploded(files)
    }
  });
}

function funcToChatImageUpload() {
  document.getElementById('chatUploadImage').addEventListener('click', function() {
    document.getElementById('chatFileUploader').click();  // Trigger a click on the hidden file input
  });

  document.getElementById('chatFileUploader').addEventListener('change', function() {
      const files = document.getElementById('chatFileUploader').files;
      console.log(files);
      handleChatUploded(files)
  });
}

function handleChatUploded(files) {
  for (const file of files) {
    if (file.type.startsWith('image/')) {
      document.getElementById('chatUplodedImageShowBack').style.display = 'flex'
      const imgElement = document.getElementById("showChatUploadedImage");
      imgElement.src = URL.createObjectURL(file);
      scrollDown()
    } 
  }
}

function sendMessage() {
  var type = false
  var msg = document.getElementById("chat-input").value;
  if (msg == "") {
    setTimeout(function() {
      document.getElementById("chat-input").value = ""
      document.getElementById("chat-input").focus();
    }, 0);
    console.log("SEND MESSAGE", document.getElementById('chatFileUploader'));
    if (document.getElementById('chatFileUploader').files.length != 0) {
      type = true
    } else {
      return
    }
  }
  var withName = document.getElementById('user-messages-username').innerHTML
  // --------------------------- IF IMAGE -------------------
  if (type) {
    const form = document.getElementById("formSkrepkaBack");
    const formData = new FormData(form);
    fetch("/chat/"+withName, {
      method: "POST",
      body: formData
  })
  .then(response => {
      if (!response.ok) {
          // throw new Error("Network response was not ok");
      }
      return response.text();
  })
  .then(data => {
    document.getElementById('messagesChatWindow').innerHTML = dataConverterChat(data, "messages-chat-window", "messagesChatWindow") 
    let msgChatScrollHeight = 0
    const messagesChatWindow1 = document.getElementById('messagesChatWindow')
    messagesChatWindow1.scrollTop = messagesChatWindow1.scrollHeight;
    msgChatScrollHeight = messagesChatWindow1.scrollTop
    document.getElementById("chat-input").focus();
      setTimeout(function() {
        document.getElementById("chat-input").value = ""
        document.getElementById("chat-input").focus();
        refreshChatOnce()
      }, 0);
      document.getElementById('chatFileUploader').value = ''
  })
  .catch(error => {
      console.error("There was a problem with the fetch operation:", error);
  });
  } else {
    fetch("/chat/"+withName, {
      method: "POST",
      headers: {
          "Content-Type": "application/x-www-form-urlencoded"
      },
      body: "message=" + encodeURIComponent(msg)
  })
  .then(response => {
      if (!response.ok) {
          // throw new Error("Network response was not ok");
      }
      return response.text();
  })
  .then(data => {
    document.getElementById('messagesChatWindow').innerHTML = dataConverterChat(data, "messages-chat-window", "messagesChatWindow") 
    let msgChatScrollHeight = 0
    const messagesChatWindow1 = document.getElementById('messagesChatWindow')
    messagesChatWindow1.scrollTop = messagesChatWindow1.scrollHeight;
    msgChatScrollHeight = messagesChatWindow1.scrollTop
    document.getElementById("chat-input").focus();
      setTimeout(function() {
        document.getElementById("chat-input").value = ""
        document.getElementById("chat-input").focus();
        refreshChatOnce()
      }, 0);
  })
  .catch(error => {
      console.error("There was a problem with the fetch operation:", error);
  });
  }
}

// -------------------- Category Sort Function ------------------- //

function sortCategory(category) {
  fetch("/post/category/"+category, {
      method: "POST",
      headers: {
          "Content-Type": "application/x-www-form-urlencoded"
      },
  })
  .then(response => {
      if (!response.ok) {
          throw new Error("Network response was not ok");
      }
      return response.text();
  })
  .then(data => {
    document.getElementById("post-table").innerHTML = dataConverter(data, "post_table", "post-table")
  })
  .catch(error => {
      console.error("There was a problem with the fetch operation:", error);
  });
}

// -------------------- Post Reaction Request ------------------- //

function reactionPost(request) {
  fetch(request, {
      method: "GET",
      headers: {
          "Content-Type": "application/x-www-form-urlencoded"
      },
  })
  .then(response => {
      if (!response.ok) {
          throw new Error("Network response was not ok");
      }
      return response.text();
  })
  .then(data => {
    document.querySelector('.post-box-reaction-back').innerHTML = dataConverter(data, "post-box-reaction-back")
  })
  .catch(error => {
      console.error("There was a problem with the fetch operation:", error);
  });
}

// -------------------- Comment Reaction Request ------------------- //

function reactionComment(request) {
  fetch(request, {
      method: "GET",
      headers: {
          "Content-Type": "application/x-www-form-urlencoded"
      },
  })
  .then(response => {
      if (!response.ok) {
          throw new Error("Network response was not ok");
      }
      return response.text();
  })
  .then(data => {
    document.querySelector('.comments').innerHTML = dataConverter(data, "comments")
  })
  .catch(error => {
      console.error("There was a problem with the fetch operation:", error);
  });
}

// -------------------- Create Comment Request ------------------- //

function createCommentFunc(request) {
  fetch(request, {
      method: "POST",
      headers: {
          "Content-Type": "application/x-www-form-urlencoded"
      },
      body: "comment=" + encodeURIComponent(document.querySelector('.comment-content').value)
  })
  .then(response => {
      if (!response.ok) {
          throw new Error("Network response was not ok");
      }
      return response.text();
  })
  .then(data => {
    document.querySelector('.comments').innerHTML = dataConverter(data, "comments")
    document.querySelector('.comment-content').value = ""
  })
  .catch(error => {
      console.error("There was a problem with the fetch operation:", error);
  });
}

// -------------------- Profile Change ------------------- //

const changeBtn = document.getElementById('change_btn')
const changeProfilePhoto = document.getElementById('change-profile-photo')
const changeBtnCancel = document.getElementById('change_btn_cancel')

function openImageFieldProfile() {
  changeProfilePhoto.style.display = 'block'
  changeBtn.style.display = 'none'
}
function closeImageFieldProfile() {
  changeProfilePhoto.style.display = 'none'
  changeBtn.style.display = 'block'
}

function changeProfilePhotoFunc(request) {
  const form = document.getElementById("myForm");
  const formData = new FormData(form);

  fetch(request, {
    method: "POST",
    body: formData,
  })
    .then(response => {
      if (!response.ok) {
        throw new Error("Network response was not ok");
    }
    return response.text();
    })
    .then((data) => {
      document.querySelector('.profile-user-in').innerHTML = dataConverter(data, "profile-user-in")
      document.getElementById('change_btn').addEventListener('click', openImageFieldProfile);
    })
    .catch((error) => {
      console.error("Error:", error);
    });
}

// -------------------- Activity Item Title Changer ------------------- //

var activityItemTypes = document.querySelectorAll('#change-type');

// Loop through each element
activityItemTypes.forEach(function(activityItemType) {
  // Apply switch case logic to each element
  switch (activityItemType.innerHTML) {
    case "createpost":
      activityItemType.innerHTML = "Сіз пост құрдыңыз";
      break;
    case "createcomment":
      activityItemType.innerHTML = "Сіз комментарий жаздыңыз";
      break;
    case "reactionpost":
      activityItemType.innerHTML = "Пост реакциясы";
      break;
    case "reactioncomment":
      activityItemType.innerHTML = "Комментарий реакциясы";
      break;
    default:
      activityItemType.innerHTML = "";
  }
});

// -------------------- Users Switcher -------------------- 

function runSwitcher() {
  if (document.getElementById('all-users-btn') != null && document.getElementById('my-chats-btn') != null) {
    document.getElementById('all-users-btn').addEventListener('click', () => {
      document.getElementById('all-exist-chats-container').classList.add('active')
      document.getElementById('all-users-btn').classList.add('font-active')
      document.getElementById('my-chats-btn').classList.remove('font-active')
      document.getElementById('all-exist-chats-container').scrollTop = 0
    });
    
    document.getElementById('my-chats-btn').addEventListener('click', () => {
      document.getElementById('all-exist-chats-container').classList.remove('active')
      document.getElementById('my-chats-btn').classList.add('font-active')
      document.getElementById('all-users-btn').classList.remove('font-active')
      document.getElementById('all-exist-chats-container').scrollTop = 0
    });
  }
}
runSwitcher()

function activateChatBackBtn() {
  document.getElementById("chat-arrow-back-btn").addEventListener('click', ()=> {
    document.getElementById('my-chats-box').style.display = 'flex'
    document.getElementById('chat-message-box').style.display = 'none'
  })
}

if (document.getElementById("chat-arrow-back-btn")) {
  activateChatBackBtn()
}

// ----------------- Search Function ------------------//

const myChatsBoxDiv = document.getElementById('myChatsBox');
const userRows = Array.from(myChatsBoxDiv.getElementsByClassName('user-row'));

const allUserBoxDiv = document.getElementById('all-user-box');
const allUserRows = Array.from(allUserBoxDiv.getElementsByClassName('user-row'));

function chatsSearchFunction() {
  filterChats(document.getElementById('chatsSearch').value)
}

    function filterChats(filterString) {
        myChatsBoxDiv.innerHTML = '';
            const filteredRows1 = userRows.filter(row => {
            const chatWith1 = row.querySelector('.chat-with').innerText.toLowerCase();
            return chatWith1.toLowerCase().includes(filterString.toLowerCase());
        });
        filteredRows1.forEach(row => {
            myChatsBoxDiv.appendChild(row.parentNode);
        });
        allUserBoxDiv.innerHTML = '';
            const filteredRows2 = allUserRows.filter(row => {
            const chatWith2 = row.querySelector('.users-name').innerText.toLowerCase();
            return chatWith2.toLowerCase().includes(filterString.toLowerCase());
        });
        filteredRows2.forEach(row => {
            allUserBoxDiv.appendChild(row.parentNode);
        });
    }



// --------------------- Message Header Menu ------------------

// const messagesHeaderMenuBtn = document.getElementById('messages-chat-header-menu')
// const myChatsShow = document.getElementById('my-chats-box-show')
// const myChatsSearch = document.getElementById('messages-chat-header-search-back')

// if (messagesHeaderMenuBtn) {
//   messagesHeaderMenuBtn.addEventListener('click', () => {
//     if (myChatsShow.style.display == "none") {
//       myChatsBox.style.width = "30%"
//       chatMessageBox.style.width = "70%"
//       myChatsShow.style.display = "block"
//       myChatsSearch.style.display = "flex"
//     } else {
//       myChatsBox.style.width = "10%"
//       chatMessageBox.style.width = "100%"
//       myChatsShow.style.display = "none"
//       myChatsSearch.style.display = "none"
//     }
//   }) 
// }