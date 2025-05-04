    var fileInput = document.getElementById("promoImage");
    var previewImg = document.getElementById("promoImagePreview");
    
    fileInput.addEventListener("change", function (event) {
        const file = event.target.files[0];
        if (file) {
            const previewUrl = URL.createObjectURL(file);
            previewImg.src = previewUrl;
        }
    });
