//WRAPPER MAIN ELEMENTS
const box = document.querySelector('.box')
const create = document.querySelector('.create')
const post = document.querySelector('.post_main')
const wrapper = document.querySelector('.wrapper')
const loginLink = document.querySelector('.login-link')
const registerLink = document.querySelector('.register-link')

//ADAPTIVE MENU ELEMENTS
const menu = document.querySelector(".icon-menu")
const navLink = document.querySelector(".navigation")
const btnPop = document.querySelector('.btnLogin')

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
  if (wrapper != null) {
    wrapper.classList.remove('active-pop')
  } else if (box != null) {
    box.classList.remove('active-pop')
  } else if (create != null) {
    create.classList.remove('active-pop')
  } else if (post != null) {
    post.classList.remove('active-pop')
  } 
  navLink.classList.toggle('mob-menu')
})

if (registerLink != null && loginLink != null) {
  registerLink.addEventListener('click', () => {
    wrapper.classList.add('active')
    wrapper.classList.remove('error')
    error_window_box.style.display = "none"
    error_window_box2.style.display = "none"
  });
  
  loginLink.addEventListener('click', () => {
    wrapper.classList.remove('active')
    wrapper.classList.remove('error')
    error_window_box.style.display = "none"
    error_window_box2.style.display = "none"
  });
}

// -------------------- Check Login Parametrs ------------------- //
//PASSWORD LOCK ICONS 
const logIconLock1 = document.querySelector('.log-icon-lock1')
const logIconLock2 = document.querySelector('.log-icon-lock2')
const logIconLock3 = document.querySelector('.log-icon-lock3')

const locks = [0, 0, 0];

function toggleLock(lockIndex, iconElement, passwordElement) {
  const lock = locks[lockIndex];
  if (lock % 2 === 0) {
    iconElement.setAttribute('name', 'lock-open');
    passwordElement.setAttribute('type', 'text');
  } else {
    iconElement.setAttribute('name', 'lock-closed');
    passwordElement.setAttribute('type', 'password');
  }
  locks[lockIndex]++;
}

//PASSWORD LABELS
const password_lbl_log = document.getElementById('password-label-log')
const password_lbl_reg = document.getElementById('password-label-reg')
const password_lbl_reg2 = document.getElementById('password-label-reg2')

if (logIconLock1 && logIconLock2 && logIconLock3) {
  logIconLock1.addEventListener('click', () => toggleLock(0, logIconLock1, password_lbl_log));
  logIconLock2.addEventListener('click', () => toggleLock(1, logIconLock2, password_lbl_reg));
  logIconLock3.addEventListener('click', () => toggleLock(2, logIconLock3, password_lbl_reg2));
}

// -------------------- Check Login Parametrs ------------------- //

//INPUTS
const logEmail = document.querySelector('.log-email')
const logPassword = document.querySelector('.log-pass')
const regUsername = document.querySelector('.reg-username')
const regEmail = document.querySelector('.reg-email')
const regPassword = document.querySelector('.reg-password')
const regPasswordSec = document.querySelector('.reg-password-second')

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

loginErrors()

function loginErrors() {
  if (form_act != null) {
    var form_act_inner = form_act.innerHTML.toString()
    if (form_act_inner === 'Log') {
      wrapper.classList.add('active-pop')
      wrapper.classList.remove('active')
      if (form_error.innerHTML.toString() !== "") {
        error_window_box.style.display = "block"
        error_window_log.innerHTML = form_error.innerHTML.toString()
        wrapper.classList.add('error')
        //SET VALUES FROM BACK
        logEmail.value = ReturnedEmail.innerHTML
        console.log("awd" + ReturnedEmail.innerHTML.toString());
        logPassword.value = ReturnedPass.innerHTML
      } else {
        wrapper.classList.remove('error')
        error_window_box.style.display = "none"
        error_window_box2.style.display = "none"
      }
    } else if (form_act_inner === 'Reg') {
      wrapper.classList.add('active-pop')
      wrapper.classList.add('active')
      if (form_error.innerHTML.toString() !== "") {
        error_window_box2.style.display = "block"
        error_window_log2.innerHTML = form_error.innerHTML.toString()
        wrapper.classList.add('error')
        //SET VALUES FROM BACK
        regUsername.value = ReturnedUsername.innerHTML
        regEmail.value = ReturnedEmail.innerHTML
        regPassword.value = ReturnedPass.innerHTML
        regPasswordSec.value = ReturnedPass2.innerHTML
      } else {
        wrapper.classList.remove('error')
        error_window_box.style.display = "none"
        error_window_box2.style.display = "none"
      }
    }
  }
}

// -------------------- Image upload Drop Down Area ------------------- //
const dropArea = document.getElementById('drop-area');
const fileInput = document.getElementById('fileInput');
const imgMessage = document.querySelector('.post-image-message')

if (dropArea != null) {
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
        fileInput.files = files
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
          imgMessage.innerHTML = "Image upload success"
          imgMessage.style.display = "block"
        } else {
          imgMessage.innerHTML = "Invalid file type!"
          imgMessage.style.display = "block"
        }
    }
  }
}

const ReturnedPostTitle = document.getElementById('ReturnedPostTitle')
const ReturnedPostContent = document.getElementById('ReturnedPostContent')
const PostTitle = document.getElementById('post-title')
const PostContent = document.getElementById('post-content')

if (ReturnedPostTitle != null) {
  PostTitle.value = ReturnedPostTitle.innerHTML
  PostContent.value = ReturnedPostContent.innerHTML
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

// -------------------- CHAT ------------------- //
var inputElementChat = document.getElementById("chat-input");
if (inputElementChat != null) {
  window.onload = function() {
    scrollToBottom();
    inputElementChat.focus();
  };
}

// -------------------- Chat Auto Scroll Down ------------------- //
let msgChatScrollHeight = 0

function scrollToBottom() {
  var messagesChatWindow = document.getElementById('messagesChatWindow');
  if (messagesChatWindow) {
    messagesChatWindow.scrollTop = messagesChatWindow.scrollHeight;
  }
  msgChatScrollHeight = msgChatWindow.scrollTop
}

// -------------------- Chat Enter Press Event ------------------- //

function handleKeyDown(event) {
  if (event.keyCode === 13) {
    sendMessage();
    inputElementChat.focus();
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

var msgChatWindow = document.getElementById('messagesChatWindow')
if (msgChatWindow != null) {
  refreshChat()
} 

let lastScrollTop = 0;
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
if (msgChatWindow != null) {
  msgChatWindow.addEventListener("scroll", debounce(function() {
    let scrollPosition = msgChatWindow.scrollTop;
      
    if (msgChatScrollHeight - scrollPosition > 20) {
      scrollDownBtn.style.display = "block";
      scrollDownBtn.style.animation = 'dropDownBtnIn 0.5s forwards';
    } else {
      scrollDownBtn.style.animation = 'dropDownBtnOut 0.5s forwards';
      setTimeout(() => {
        scrollDownBtn.style.display = "none";
      }, 500);
    }
  }, 20));
}

if (scrollDownBtn != null) {
  scrollDownBtn.addEventListener('click', () => {
    scrollToBottom()
  })
}

function refreshChat() {
  setInterval(() => {
    var withName = document.getElementById('user-messages-username').innerHTML
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
      var oldChat = document.getElementById("messagesChatWindow").innerHTML
      if (oldChat != dataConverter(data, "messages-chat-window", "messagesChatWindow")) {
        document.getElementById("messagesChatWindow").innerHTML = dataConverter(data, "messages-chat-window", "messagesChatWindow")
        scrollToBottom();
      }
    })
    .catch(error => {
        console.error("There was a problem with the fetch operation:", error);
    });
  }, 1000)
}

// -------------------- Chat Send Message Func ------------------- //

function sendMessage() {
  var msg = document.getElementById("chat-input").value;
  if (msg == "") {
    setTimeout(function() {
      document.getElementById("chat-input").value = ""
      document.getElementById("chat-input").focus();
    }, 0);
    return
  }
  var withName = document.getElementById('user-messages-username').innerHTML
  fetch("/chat/"+withName, {
      method: "POST",
      headers: {
          "Content-Type": "application/x-www-form-urlencoded"
      },
      body: "message=" + encodeURIComponent(msg)
  })
  .then(response => {
      if (!response.ok) {
          throw new Error("Network response was not ok");
      }
      return response.text();
  })
  .then(data => {
    document.getElementById("messagesChatWindow").innerHTML = dataConverter(data, "messages-chat-window", "messagesChatWindow") 
      scrollToBottom();
      setTimeout(function() {
        document.getElementById("chat-input").value = ""
        document.getElementById("chat-input").focus();
      }, 0);
  })
  .catch(error => {
      console.error("There was a problem with the fetch operation:", error);
  });
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
  console.log("s");
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
      console.log(data);
      document.getElementById('change_btn').addEventListener('click', openImageFieldProfile);
    })
    .catch((error) => {
      console.error("Error:", error);
    });
}

// function changeProfilePhoto(request) {
//   fetch(request, {
//     method: "POST",
//     headers: {
//         "Content-Type": "application/x-www-form-urlencoded"
//     },
//     body: "comment=" + encodeURIComponent(document.querySelector('.comment-content').value)
// })
// .then(response => {
//     if (!response.ok) {
//         throw new Error("Network response was not ok");
//     }
//     return response.text();
// })
// .then(data => {
//   document.querySelector('.comments').innerHTML = dataConverter(data, "comments")
//   document.querySelector('.comment-content').value = ""
// })
// .catch(error => {
//     console.error("There was a problem with the fetch operation:", error);
// });
// }

// -------------------- Activity Item Title Changer ------------------- //

var activityItemTypes = document.querySelectorAll('#change-type');

// Loop through each element
activityItemTypes.forEach(function(activityItemType) {
  // Apply switch case logic to each element
  switch (activityItemType.innerHTML) {
    case "createpost":
      activityItemType.innerHTML = "You created post";
      break;
    case "createcomment":
      activityItemType.innerHTML = "You created comment";
      break;
    case "reactionpost":
      activityItemType.innerHTML = "Post reaction";
      break;
    case "reactioncomment":
      activityItemType.innerHTML = "Comment reaction";
      break;
    default:
      activityItemType.innerHTML = "";
  }
});