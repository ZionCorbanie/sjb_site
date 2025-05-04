document.addEventListener("htmx:confirm", function(e) {
    if (!e.detail.question) return

    e.preventDefault()

    Swal.fire({
        title: "Doorgaan?",
        text: e.detail.question,
        confirmButtonText: "Verwijderen",
        confirmButtonColor: "#d33",
        showCancelButton: true,
        cancelButtonText: "Annuleren",
        cancelButtonColor: "#3085d6",
    }).then(function(result) {
      if (result.isConfirmed) {
        e.detail.issueRequest(true);
      }
    })
  })
