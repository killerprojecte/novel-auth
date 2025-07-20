import { Api } from './api';

function showForm(evt, formId) {
  const tabButtons = document.getElementsByClassName('tablinks');
  const formLogin = document.getElementById('form-login');
  const formRegister = document.getElementById('form-register');

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
  .addEventListener('click', (event) => showForm(event, 'login'));

document
  .getElementById('tab-register')
  .addEventListener('click', (event) => showForm(event, 'register'));

function getQueryParam(paramName, defaultValue) {
  const params = new URLSearchParams(window.location.href);
  return params.get(paramName) || defaultValue;
}

function redirectAfterLogin() {
  const defaultRedirect = '/';
  const redirectUrl = getQueryParam('redirect', defaultRedirect);
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
  const app = getQueryParam('from', '');
  const username = document.getElementById('username').value;
  const password = document.getElementById('password').value;

  if (Api.login.isPending) return;
  Api.login(app, username, password)
    .then(() => redirectAfterLogin())
    .catch((error) => alert(`登录失败: ${error}`));
});

document.getElementById('btn-register').addEventListener('click', (event) => {
  event.preventDefault();
  const app = getQueryParam('from', '');
  const username = document.getElementById('r-username').value;
  const password = document.getElementById('r-password').value;
  const email = document.getElementById('r-email').value;
  const verifyCode = document.getElementById('r-verify-code').value;

  if (Api.register.isPending) return;
  Api.register(app, email, username, password, verifyCode)
    .then(() => redirectAfterLogin())
    .catch((error) => alert(`注册失败: ${error}`));
});

document
  .getElementById('btn-send-verify-code')
  .addEventListener('click', (event) => {
    event.preventDefault();
    const email = document.getElementById('r-email').value;

    if (Api.requestVerifyCode.isPending) return;
    Api.requestVerifyCode(email)
      .then(() => alert('验证码已发送到您的邮箱'))
      .catch((error) => alert(`发送验证码失败: ${error}`));
  });
