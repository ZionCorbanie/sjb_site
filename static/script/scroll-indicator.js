document.addEventListener("DOMContentLoaded", function () {
    const indicator = document.getElementById("scroll-indicator");

    window.addEventListener("scroll", function () {
        indicator.classList.add("opacity-0");
        setTimeout(() => indicator.classList.add("hidden"), 500);
    });
});

function scrollToTarget() {
    let yOffset = -75;
    let y = document.getElementById("posts").getBoundingClientRect().top + window.scrollY + yOffset;
    window.scrollTo({top: y, behavior: 'smooth'});

    const indicator = document.getElementById("scroll-indicator");
    indicator.classList.add("opacity-0");
    setTimeout(() => indicator.classList.add("hidden"), 500);
}
