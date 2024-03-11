//WRAPPER MAIN ELEMENTS
const box = document.querySelector('.box')
const create = document.querySelector('.create')
const post = document.querySelector('.post_main')
const wrapper = document.querySelector('.wrapper')
const loginLink = document.querySelector('.login-link')
const registerLink = document.querySelector('.register-link')
const btnPop = document.querySelector('.btnLogin')

//ADAPTIVE MENU ELEMENTS
const menu = document.querySelector(".icon-menu")
const navLink = document.querySelector(".navigation")

//PASSWORD LABELS
const password_lbl_log = document.getElementById('password-label-log')
const password_lbl_reg = document.getElementById('password-label-reg')
const password_lbl_reg2 = document.getElementById('password-label-reg2')

//PASSWORD LOCK ICONS 
const logIconLock1 = document.querySelector('.log-icon-lock1')
const logIconLock2 = document.querySelector('.log-icon-lock2')
const logIconLock3 = document.querySelector('.log-icon-lock3')

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

var headerDiv = document.querySelector("header");

window.addEventListener("scroll", function () {
  if (window.scrollY >= 15) {
    headerDiv.style.display = "none"
  }
  else if (window.scrollY < 15) {
    headerDiv.style.display = "flex"
  }
});

var lock1 = 0
var lock2 = 0
var lock3 = 0

if (logIconLock1 != null && logIconLock2 != null && logIconLock3 != null) {
  logIconLock1.addEventListener('click', function () {
    if (lock1 % 2 == 0) {
      logIconLock1.setAttribute('name', 'lock-open')
      password_lbl_log.setAttribute('type', 'text')
    } else {
      logIconLock1.setAttribute('name', 'lock-closed')
      password_lbl_log.setAttribute('type', 'password')
    }
    lock1++
  })
  
  logIconLock2.addEventListener('click', function () {
    if (lock2 % 2 == 0) {
      logIconLock2.setAttribute('name', 'lock-open')
      password_lbl_reg.setAttribute('type', 'text')
    } else {
      logIconLock2.setAttribute('name', 'lock-closed')
      password_lbl_reg.setAttribute('type', 'password')
    }
    lock2++
  })
  
  logIconLock3.addEventListener('click', function () {
    if (lock3 % 2 == 0) {
      logIconLock3.setAttribute('name', 'lock-open')
      password_lbl_reg2.setAttribute('type', 'text')
    } else {
      logIconLock3.setAttribute('name', 'lock-closed')
      password_lbl_reg2.setAttribute('type', 'password')
    }
    lock3++
  })
}

//GET FROM BACK INPUT VALUES
const ReturnedEmail = document.getElementById('ReturnedEmail')
const ReturnedUsername = document.getElementById('ReturnedUsername')
const ReturnedPass = document.getElementById('ReturnedPass')
const ReturnedPass2 = document.getElementById('ReturnedPass2')

//TO GET ERRORS FROM BACK
const form_act = document.getElementById("form-act")
const form_error = document.getElementById("form-error")

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

const dropArea = document.getElementById('drop-area');
const fileInput = document.getElementById('fileInput');
const imgMessage = document.querySelector('.post-image-message')

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
