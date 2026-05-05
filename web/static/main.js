// [SUA] JSON-only chat rendering (khong dung HTMX).
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

    function appendUserMessage(message) { // [THEM_MOI]
      var article = document.createElement("article");
      article.className = "msg user";
      var h4 = document.createElement("h4");
      h4.textContent = "Ban";
      var p = document.createElement("p");
      p.textContent = message;
      article.appendChild(h4);
      article.appendChild(p);
      result.appendChild(article);
    }

    function appendBotMessage(data) { // [THEM_MOI]
      var article = document.createElement("article");
      article.className = "msg";
      if (data && data.is_order_completed) {
        article.classList.add("ok");
      }

      var h4 = document.createElement("h4");
      h4.textContent = "Phan hoi";
      article.appendChild(h4);

      if (data && data.text) {
        var pText = document.createElement("p");
        pText.textContent = data.text;
        article.appendChild(pText);
      }

      if (data && data.is_order_completed) {
        var pOrder = document.createElement("p");
        pOrder.textContent = "Ma don: " + (data.order_id || "");
        article.appendChild(pOrder);
        if (data.payment_qr_url) {
          var pQR = document.createElement("p");
          pQR.textContent = "Quet ma QR de thanh toan:";
          article.appendChild(pQR);

          var img = document.createElement("img");
          img.src = data.payment_qr_url;
          img.alt = "QR thanh toan don hang " + (data.order_id || "");
          img.loading = "lazy";
          img.style.maxWidth = "240px";
          img.style.height = "auto";
          img.style.border = "1px solid #eee";
          img.style.borderRadius = "8px";
          img.style.padding = "8px";
          img.style.background = "#fff";
          article.appendChild(img);
        }
      }

      if (data && Array.isArray(data.photourl) && data.photourl.length > 0) {
        var photos = document.createElement("div");
        photos.className = "media-grid";
        data.photourl.forEach(function (url) {
          if (!url) return;
          var figure = document.createElement("figure");
          figure.className = "media-item";
          var img = document.createElement("img");
          img.src = url;
          img.alt = "Hinh anh tu agent";
          img.loading = "lazy";
          img.style.maxWidth = "100%";
          img.style.height = "auto";
          img.style.borderRadius = "8px";
          img.style.border = "1px solid #eee";
          figure.appendChild(img);
          photos.appendChild(figure);
        });
        article.appendChild(photos);
      }

      if (data && Array.isArray(data.videourl) && data.videourl.length > 0) {
        var videos = document.createElement("div");
        videos.className = "media-grid";
        data.videourl.forEach(function (url) {
          if (!url) return;
          var figure = document.createElement("figure");
          figure.className = "media-item";
          var video = document.createElement("video");
          video.src = url;
          video.controls = true;
          video.preload = "metadata";
          video.style.maxWidth = "100%";
          video.style.height = "auto";
          video.style.borderRadius = "8px";
          video.style.border = "1px solid #eee";
          figure.appendChild(video);
          videos.appendChild(figure);
        });
        article.appendChild(videos);
      }

      result.appendChild(article);
    }

    form.addEventListener("submit", async function (event) {
      event.preventDefault();

      var messageInput = form.querySelector('input[name="message"]');
      if (!messageInput) return;
      var message = (messageInput.value || "").trim();
      if (!message) return;

      appendUserMessage(message); // [THEM_MOI]

      setLoading(true);
      try {
        var formData = new FormData();
        formData.set("message", message);

        var response = await fetch(form.action || "/chat", {
          method: "POST",
          body: formData
        });

        if (!response.ok) {
          throw new Error("Request failed with status " + response.status);
        }

        var data = await response.json(); // [SUA]
        appendBotMessage(data); // [THEM_MOI]
        form.reset();
      } catch (error) {
        var safeMsg = (error && error.message) ? error.message : "Unknown error";
        result.insertAdjacentHTML(
            "beforeend",
          '<article class="msg err"><h4>Loi</h4><p>' + safeMsg + "</p></article>"
        );
      } finally {
        result.scrollTop = result.scrollHeight; // [THEM_MOI]
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
