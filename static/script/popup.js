document.addEventListener('click', function(event) {
    let popup = document.getElementById("popup");
    let popupWindow = document.getElementById("popupContent");
    let button = document.getElementById("closePopup");

    if (popupWindow != null && !popupWindow.contains(event.target)) {
        popup.remove();
    }
    if (button != null && button.contains(event.target)) {
        popup.remove();
    }
});
