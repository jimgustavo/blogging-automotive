//  static/script.js

async function getBlogPosts() {
    try {
        const response = await fetch('http://localhost:8080/blogposts');
        if (!response.ok) {
            throw new Error('Failed to fetch blog posts');
        }
        
        const data = await response.json();
        
        // Check if the response contains data
        if (data && data.length > 0) {
            console.log('Blog posts:');
            console.log(data);
            
            // Get the container element where you want to display the blog posts
            const container = document.getElementById('blog-cards');
            
            // Clear any existing content in the container
            container.innerHTML = '';
            
            // Iterate through the data and generate HTML for each blog post
            data.forEach(post => {
                const blogCard = document.createElement('div');
                blogCard.classList.add('blog-card');
                
                // Generate HTML for the blog post
                blogCard.innerHTML = `
                        <img src=${post.picture} alt="Blog 4">
                        <div class="blog-text">
                            <h3>${post.title}</h3>
                            <p>${post.summary}...<a href="/static/reading-page.html?id=${ post.id }">Read more</a></p>
                            <!-- Generate link with the post ID as part of the URL -->
                            
                            <p>Author: ${post.author}</p>
                            <p>Created At: ${post.created_at}</p>
                        </div>
                `;
                //<a href={'/post/'+item._id} style={{cursor:'pointer', color: 'gray'}} key={i} className={classes.blog}>
                
                // Append the blog post element to the container
                container.appendChild(blogCard);
            });
        } else {
            console.log('No blog posts available');
        }
    } catch (error) {
        console.error('Error fetching blog posts:', error);
    }
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


// Fetch and display auto parts when the page loads
window.onload = getBlogPosts;
