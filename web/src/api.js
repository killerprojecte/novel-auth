async function post(url, body) {
  const controller = new AbortController();
  const signal = controller.signal;

  // 设置5秒超时
  const timeoutId = setTimeout(() => controller.abort(), 5000);

  try {
    const response = await fetch('/api/v1/auth/' + url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
      signal,
    });
    clearTimeout(timeoutId);
    if (!response.ok) {
      const message = await response.text();
      console.log(response.statusText);
      throw `[${response.status}] ${message || response.statusText}`;
    }
    return await response.text();
  } catch (error) {
    if (error.name === 'AbortError') {
      throw '请求超时，请稍后再试';
    }
    throw `${error}`;
  }
}

function debounce(func) {
  const newFunc = async function (...args) {
    if (newFunc.isPending) return;
    newFunc.isPending = true;
    try {
      return await func(...args);
    } finally {
      newFunc.isPending = false;
    }
  };
  newFunc.isPending = false;
  return newFunc;
}

export const Api = {
  register: debounce((app, email, username, password, verifyCode) =>
    post('register', {
      app,
      email,
      username,
      password,
      'verify-code': verifyCode,
    })
  ),
  login: debounce((app, username, password) =>
    post('login', {
      app,
      username,
      password,
    })
  ),
  requestVerifyCode: debounce((email) =>
    post('request-verify-code', {
      email,
    })
  ),
};
