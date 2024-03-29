const form = document.querySelector("#fwt")
form.addEventListener('htmx:afterRequest', function(event) {
  const messageBox = document.createElement('div');
  if (event.detail.successful) {
    messageBox.textContent = 'OK ✅';
    messageBox.style.color = 'green';
  } else {
    messageBox.textContent = 'Failed ❌';
    messageBox.style.color = 'red';
  }
  const requestInitiator = event.target;
  requestInitiator.appendChild(messageBox);
  setTimeout(() => {
    messageBox.remove();
  }, 3000);
});

form.addEventListener("htmx:afterRequest", (e) => {
  e.preventDefault()
  let ipts = e.currentTarget.querySelectorAll("input[type='text']")
  ipts.forEach(element => {
    element.value = ""
  });
})
