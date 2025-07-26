async function post(url: string, body: any) {
  const controller = new AbortController();
  const signal = controller.signal;

  const delay = () => new Promise((resolve) => setTimeout(resolve, 5000));
  await delay(); // 延迟1秒，模拟网络延迟

  // 设置5秒超时
  const timeoutId = setTimeout(() => controller.abort(), 5000);

  try {
    const response = await fetch("/api/v1/auth/" + url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
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
    if (error instanceof Error && error.name === "AbortError") {
      throw "请求超时，请稍后再试";
    }
    throw `${error}`;
  }
}

function debounce<T extends (...args: any[]) => Promise<any>>(func: T) {
  const newFunc = async function (
    ...args: Parameters<T>
  ): Promise<ReturnType<T> | undefined> {
    if (newFunc.isPending) return undefined;
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
  register: debounce((app, username, password, email, otp) =>
    post("register", {
      app,
      username,
      password,
      email,
      otp,
    }),
  ),
  login: debounce((app, username, password) =>
    post("login", {
      app,
      username,
      password,
    }),
  ),
  requestOtp: debounce((email) =>
    post("request-otp", {
      email,
    }),
  ),
  resetPassword: debounce((email, password, otp) =>
    post("reset-password", {
      email,
      otp,
      password,
    }),
  ),
};

function getQueryParam(paramName: string, defaultValue: string) {
  const params = new URLSearchParams(window.location.href);
  return params.get(paramName) || defaultValue;
}

export function redirectAfterLogin() {
  const defaultRedirect = "/";
  const redirectUrl = getQueryParam("redirect", defaultRedirect);
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
