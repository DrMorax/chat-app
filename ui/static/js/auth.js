export const authenticator = {
  authenticate: function () {
    let data = { email: "test4@test.com", password: "12345678" };

    fetch("/user/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    }).then((res) => {
      console.log("Request complete! response:", res);
    });
  },
  register: function () {
    let data = { name: "test4", email: "test4@test.com", password: "12345678" };

    fetch("/user/signup", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    }).then((res) => {
      console.log("Request complete! response:", res);
    });
  },
};
