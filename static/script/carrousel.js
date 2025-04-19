document.addEventListener("DOMContentLoaded", function () {
    const carousel = document.getElementById("carousel");
    const images = carousel.children;
    const prevBtn = document.getElementById("prevBtn");
    const nextBtn = document.getElementById("nextBtn");

    let index = 0;
    const totalImages = images.length;
    const intervalTime = 5000; // 5 seconds

    function updateCarousel() {
        carousel.style.transform = `translateX(-${index * 100}%)`;
    }

    function nextSlide() {
        index = (index + 1) % totalImages;
        updateCarousel();
    }

    function prevSlide() {
        index = (index - 1 + totalImages) % totalImages;
        updateCarousel();
    }

    // Auto-scroll every 5 seconds
    // let autoScroll = setInterval(nextSlide, intervalTime);

    // Pause auto-scroll on button click, then restart
    function resetAutoScroll() {
        clearInterval(autoScroll);
        autoScroll = setInterval(nextSlide, intervalTime);
    }

    nextBtn.addEventListener("click", function () {
        nextSlide();
        // resetAutoScroll();
    });

    prevBtn.addEventListener("click", function () {
        prevSlide();
        // resetAutoScroll();
    });
});
