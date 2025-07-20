const tabButtons = document.getElementsByClassName('tablinks');
const formLogin = document.getElementById('form-login');
const formRegister = document.getElementById('form-register');

function openForm(evt, formId) {
  Array.from(tabButtons).forEach((btn) => btn.classList.remove('active'));
  evt.currentTarget.classList.add('active');
  if (formId == 'login') {
    formRegister.classList.remove('active');
    formLogin.classList.add('active');
  } else {
    formLogin.classList.remove('active');
    formRegister.classList.add('active');
  }
}

document
  .getElementById('tab-login')
  .addEventListener('click', (event) => openForm(event, 'login'));

document
  .getElementById('tab-register')
  .addEventListener('click', (event) => openForm(event, 'register'));

function redirectAfterLogin() {
  const defaultRedirect = '/';
  const urlParams = new URLSearchParams(window.location.search);
  const redirectUrl = urlParams.get('redirect') || defaultRedirect;
  try {
    const finalUrl = new URL(redirectUrl, window.location.origin);
    if (finalUrl.origin === window.location.origin) {
      window.location.href =
        finalUrl.pathname + finalUrl.search + finalUrl.hash;
    } else {
      window.location.href = defaultRedirect;
    }
  } catch (e) {
    window.location.href = defaultRedirect;
  }
}

document.getElementById('btn-login').addEventListener('click', (event) => {
  event.preventDefault();
  const username = document.getElementById('username').value;
  const password = document.getElementById('password').value;
  redirectAfterLogin();
});

document.getElementById('btn-register').addEventListener('click', (event) => {
  event.preventDefault();
  const username = document.getElementById('r-username').value;
  const password = document.getElementById('r-password').value;
  redirectAfterLogin();
});

document
  .getElementById('btn-send-verify-code')
  .addEventListener('click', (event) => {
    event.preventDefault();
    const email = document.getElementById('r-email').value;
  });
