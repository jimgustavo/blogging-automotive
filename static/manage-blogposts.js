//script.js

// Wait for the DOM content to be fully loaded before initializing CKEditor
document.addEventListener("DOMContentLoaded", function () {
    ClassicEditor
        .create(document.querySelector('#editor'), {
            plugins: [
                'Heading', // For titles and subtitles
                'Image', // For images
                'MediaEmbed', // For videos
                'Essentials', // For basic text formatting
                'List', // For bullet points
                'Link', // For links
                'Table', // For creating and editing tables
                // Remove SourceEditing plugin
                // Add more plugins as needed
            ],
            toolbar: ['heading', 'link', 'bulletedList', 'numberedList', 'mediaEmbed', 'insertTable'] // Define toolbar options
        })
        .then(editor => {
            console.log('Editor was initialized');
            window.editor = editor; // Store editor instance globally
        })
        .catch(error => {
            console.error('Error initializing editor', error);
        });
});

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
            const container = document.getElementById('card-container');
            
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
                            <p>Author: ${post.author}</p>
                            <p>Created At: ${post.created_at}</p>
                            <p>Category: ${post.category}</p>
                            <button onclick="updateBlogPost('${post.id}')">Update</button>
                            <button onclick="deleteBlogPost('${post.id}')">Delete</button>

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

// Function to handle the click event for creating a new blog post
async function createBlogPost() {
    try {
        const form = document.getElementById('createBlogPostForm');
        const category = form.elements['category'].value;
        const title = form.elements['title'].value;
        const picture = form.elements['picture'].value;
        const summary = form.elements['summary'].value;
        const author = form.elements['author'].value;
        const editor_data = window.editor.getData(); // Access CKEditor instance from the global window object

        const jsonData = {
            category: category, 
            title: title,
            picture: picture,
            summary: summary,
            author: author,
            editor_data: editor_data,
        };

        console.log('JSON Data:', jsonData);

        const response = await fetch('/blogposts', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(jsonData)
        });

        if (response.ok) {
            console.log('Blog post created successfully');
            getBlogPosts()
            // Optionally, you can redirect to another page or perform other actions after successful creation
        } else {
            console.error('Failed to create blog post');
            // Optionally, you can show an error message to the user
        }
    } catch (error) {
        console.error('Error creating blog post:', error);
        // Optionally, you can show an error message to the user
    }
}

// Function to handle the click event for updating a blog post
async function updateBlogPost(id) {
    try {
        // Fetch the existing blog post data
        const response = await fetch(`/blogposts/${id}`);
        if (!response.ok) {
            throw new Error('Failed to fetch blog post for update');
        }
        const postData = await response.json();

        const popup = document.createElement('div');
        popup.className = 'popup-container'; // Add class to popup container
        popup.innerHTML = `
            <form id="updateBlogPostForm" class="update-form">
                <h1>Update the Blog Post</h1>
                <label for="updateTitle">Title:</label><br>
                <input type="text" id="updateTitle" name="title" value="${postData.title}" required><br>
                <label for="updatePicture">Picture:</label><br>
                <input type="text" id="updatePicture" name="picture" value="${postData.picture}" required><br>
                <label for="updateSummary">Summary:</label><br>
                <textarea id="updateSummary" name="summary" rows="5" cols="50" required>${postData.summary}</textarea><br>
                <label for="updateEditorData">Update Editor Data:</label><br>
                <textarea id="updateEditorData" name="updateEditorData" rows="30" cols="50" required>${postData.editor_data}</textarea><br>
                <label for="updateAuthor">Author:</label><br>
                <input type="text" id="updateAuthor" name="author" value="${postData.author}" required><br>
                <label for="updateCategory">Category:</label><br>
                <input type="text" id="updateCategory" name="category" value="${postData.category}" required><br>
                <button type="submit">Update</button>
            </form>
        `;
        document.body.appendChild(popup);

        const form = document.getElementById('updateBlogPostForm');
        form.addEventListener('submit', async (event) => {
            event.preventDefault();
            const title = form.elements['title'].value;
            const picture = form.elements['picture'].value;
            const summary = form.elements['summary'].value;
            const editor_data = form.elements['updateEditorData'].value;
            const author = form.elements['author'].value;
            const category = form.elements['category'].value;

            const jsonData = {
                title: title,
                picture: picture,
                summary: summary,
                editor_data: editor_data,
                author: author,
                category: category,
            };

            const response = await fetch(`/blogposts/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(jsonData)
            });

            if (response.ok) {
                console.log('Blog post updated successfully');
                getBlogPosts();
                document.body.removeChild(popup);
            } else {
                console.error('Failed to update blog post');
                // Optionally, you can show an error message to the user
            }
        });
    } catch (error) {
        console.error('Error updating blog post:', error);
        // Optionally, you can show an error message to the user
    }
}

// Function to handle the deletion of a blog post
async function deleteBlogPost(id) {
    try {
        const confirmDelete = confirm("Are you sure you want to delete this blog post?");
        if (!confirmDelete) {
            return; // If user cancels deletion, exit the function
        }

        const response = await fetch(`/blogposts/${id}`, {
            method: 'DELETE',
        });

        if (response.ok) {
            console.log('Blog post deleted successfully');
            getBlogPosts(); // Refresh the blog post list after deletion
        } else {
            console.error('Failed to delete blog post');
            // Optionally, you can show an error message to the user
        }
    } catch (error) {
        console.error('Error deleting blog post:', error);
        // Optionally, you can show an error message to the user
    }
}

// Reading Modes
let currentMode = 0; // 0: Dark Mode, 1: Day Mode, 2: Relaxed Mode
const modes = [
    { backgroundColor: '#222', color: '#3C4C24', name: 'Dark Mode' },
    { backgroundColor: '#fff', color: '#444444', name: 'Day Mode' },
    { backgroundColor: '#FFE5B4', color: '#3C4C24', name: 'Relaxed Reading' }
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

// Set the window.onload handler to call onLoad function
window.onload = getBlogPosts;