// The Reading Page 
window.onload = function() {
    // Get the blog post ID from the URL
    const urlParams = new URLSearchParams(window.location.search);
    const postId = urlParams.get('id');

    // Function to fetch the blog post details from the API
    async function getBlogPostDetails(postId) {
        try {
            const response = await fetch(`http://localhost:8080/blogposts/${postId}`);
            if (!response.ok) {
                throw new Error('Failed to fetch blog post');
            }
            const postData = await response.json();
            return postData;
        } catch (error) {
            console.error('Error fetching blog post:', error);
        }
    }

    // Function to display the blog post details
    async function displayBlogPost() {
        const postData = await getBlogPostDetails(postId);
        if (postData) {
            document.getElementById('post-title').textContent = postData.title;
            document.getElementById('post-image').src = postData.picture;
            document.getElementById('post-content').innerHTML = postData.editor_data;
            document.getElementById('post-author').textContent = postData.author;
            document.getElementById('post-created-at').textContent = postData.created_at;
        }
    }

    displayBlogPost();
}

// Reading Modes
let currentMode = 0; // 0: Dark Mode, 1: Day Mode, 2: Relaxed Mode
const modes = [
    { backgroundColor: '#222', color: '#fff', name: 'Dark Mode' },
    { backgroundColor: '#fff', color: '#222', name: 'Day Mode' },
    { backgroundColor: '#FFE5B4', color: '#333', name: 'Relaxed Reading' }
];

const modeButton = document.getElementById('mode-button');

modeButton.addEventListener('click', function() {
    currentMode = (currentMode + 1) % modes.length;
    applyMode(currentMode);
});

function applyMode(modeIndex) {
    document.body.style.backgroundColor = modes[modeIndex].backgroundColor;
    document.body.style.color = modes[modeIndex].color;
    modeButton.textContent = modes[modeIndex].name;
}

// Apply initial mode (Dark Mode)
applyMode(currentMode);

