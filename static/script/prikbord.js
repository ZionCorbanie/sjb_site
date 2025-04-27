if (document.readyState !== 'loading') {
    prikbord();
} else {
    document.addEventListener('DOMContentLoaded', function () {
        prikbord();
    });
}

function prikbord() {
    const carousel = document.getElementById("prikbord");
    const images = carousel.children;
    const prevBtn = document.getElementById("prkPrevBtn");
    const nextBtn = document.getElementById("prkNextBtn");

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
    let autoScroll = setInterval(resetAutoScroll, intervalTime);

    // Pause auto-scroll on button click, then restart
    function resetAutoScroll() {
        clearInterval(autoScroll);
        nextSlide();
        autoScroll = setInterval(resetAutoScroll, intervalTime);
    }

    nextBtn.addEventListener("click", function () {
        nextSlide();
        clearInterval(autoScroll);
        // resetAutoScroll();
    });

    prevBtn.addEventListener("click", function () {
        prevSlide();
        clearInterval(autoScroll);
        // resetAutoScroll();
    });
}
