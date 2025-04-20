document.addEventListener("DOMContentLoaded", function () {
    const indicator = document.getElementById("scroll-indicator");

    //add flex tag and remove hidden if the page is not scrolled
    if (window.scrollY === 0) {
        indicator.classList.remove("hidden");
        indicator.classList.add("flex");
    }

    window.addEventListener("scroll", function () {
        indicator.classList.add("opacity-0");
        setTimeout(() => indicator.classList.add("hidden"), 500);
    });
});

function scrollToTarget() {
    let yOffset = -75;
    let y = document.getElementById("target").getBoundingClientRect().top + window.scrollY + yOffset;
    window.scrollTo({top: y, behavior: 'smooth'});

    const indicator = document.getElementById("scroll-indicator");
    indicator.classList.add("opacity-0");
    setTimeout(() => indicator.classList.add("hidden"), 500);
}
