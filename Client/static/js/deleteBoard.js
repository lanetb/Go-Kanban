htmx.onLoad(function (content) {
    const deleteBoardButtontoggle = document.getElementById('delete-board-button')
    const deleteBoardButtons = content.querySelectorAll('[data-board-delete-btn]')
    const boards = content.querySelectorAll('[data-board]')

    deleteBoardButtontoggle.addEventListener('click', () => {
        deleteBoardButtons.forEach(button => {
            if (button.classList.contains('delete-notactive')) {
                button.classList.remove('delete-notactive')
            } else {
                button.classList.add('delete-notactive')
            }
        });
        boards.forEach(board => {
            if (board.classList.contains('deleteable')) {
                board.classList.remove('deleteable')
            } else {
                board.classList.add('deleteable')
            }
        });
    });
});