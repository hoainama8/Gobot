// [THEM_MOI] Replace HTMX chat submit with vanilla JavaScript.
(function () {
  function setLoading(isLoading) {
    if (!document.body) return;
    if (isLoading) {
      document.body.classList.add("is-loading");
      return;
    }
    document.body.classList.remove("is-loading");
  }

  function initChatForm() {
    var form = document.getElementById("chat-form");
    var result = document.getElementById("chat-result");
    if (!form || !result) return;

    form.addEventListener("submit", async function (event) {
      event.preventDefault();

      var messageInput = form.querySelector('input[name="message"]');
      if (!messageInput) return;
      var message = (messageInput.value || "").trim();
      if (!message) return;

      setLoading(true);
      try {
        var formData = new FormData();
        formData.set("message", message);

        var response = await fetch(form.action || "/chat", {
          method: "POST",
          body: formData,
          headers: {
            "HX-Request": "true"
          }
        });

        if (!response.ok) {
          throw new Error("Request failed with status " + response.status);
        }

        var html = await response.text();
        result.insertAdjacentHTML("beforeend", html);
        form.reset();
      } catch (error) {
        var safeMsg = (error && error.message) ? error.message : "Unknown error";
        result.insertAdjacentHTML(
          "beforeend",
          '<article class="msg err"><h4>Loi</h4><p>' + safeMsg + "</p></article>"
        );
      } finally {
        setLoading(false);
      }
    });
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initChatForm);
    return;
  }
  initChatForm();
})();
