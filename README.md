# 🌐 Real Time Forum 🌐

Welcome to the **Real Time Forum**! 🚀 This repository contains the source code for a dynamic forum application built with Go and JS, following the Model-View-Controller (MVC) architectural pattern. 🛠️

## 🎯 Features 🎯

- **User Authentication**: 🔐 Users can register, log in, and log out.
- **Private Messages**: 🗨️ Users can send privates messages each other.
- **Real Time Actions**: 🔐 This forum is built as a Single Page Application, so all actions are carried out in real time.
- **Post Management**: 📝 Users can create, read, edit, and delete posts.
- **Comment Management**: 🗨️ Users can comment on posts and manage their comments.
- **Like/Dislike System**: 👍 Users can like or dislike posts and comments.
- **Session Management**: 🔄 The application manages user sessions to maintain state between requests.
- **Database Operations**: 📊 The application uses SQLite as its database and includes CRUD operations for various models.

## 🚀 Getting Started 🚀

To get started with the project, follow these steps:

1. **Clone the Repository**: Clone the repository using `git clone https://learn.zone01dakar.sn/git/edieng/real-time-forum.git`.
2. **Set Up the Database**: The application uses SQLite as its database. Ensure that SQLite is installed on your machine.
3. **Run the Application**: Start the application by running `go run .` in the project directory.
4. **Access the Application**: Open a web browser and navigate to `http://localhost:8989` to access the forum.

## 🧩 MVC Architecture 🧩

We've separated the application into three main components: Model, View, and Controller. This separation allows for efficient code organization, making the application easier to manage...

- **Model**: Defined by Go files in the `models` directory, which include data structures and database interaction methods. 
- **View**: Managed by HTML templates and CSS stylesheets in the `views` and `assets/css` directories, responsible for rendering the user interface. 
- **Controller**: Represented by Go files in the `controllers` directory, handles HTTP requests, interacts with Models, and passes data to Views. 

        ├── controllers
        │   ├── categoryController.go
        │   ├── commentController.go
        │   ├── dislikeController.go
        │   ├── likeController.go
        │   ├── mainController.go
        │   ├── postController.go
        │   ├── sessionController.go
        │   └── userController.go
        ├── models
        │   ├── category.go
        │   ├── comment.go
        │   ├── dislike.go
        │   ├── like.go
        │   ├── model.go
        │   ├── post.go
        │   ├── session.go
        │   └── user.go
        └── views
            ├── assets
            │   ├── css
            │   │   ├── auth.css
            │   │   ├── index.css
            │   │   └── style.css
            │   ├── img
            │   │   ├── bg
            │   │   │   └── forum.ico.png
            │   │   ├── posts
            │   │   │   └── stark-and-parker.png
            │   │   └── profile
            │   │       ├── default-user.png
            │   │       ├── spidey.png
            │   │       └── tony-stark.png
            │   └── js
            │       ├── auth.js
            │       ├── index.js
            │       └── main.js
            ├── auth.html
            └── index.html

This MVC architecture ensures clear responsibilities, facilitating code organization, testing, and debugging.

## 🤝 Contributing 🤝

Contributions are welcome! Please feel free to open issues for any bugs or feature requests.

## 👥 Authors 👥

- **@lodiouf**
- **@jgoudiab** 
- **@edieng**
- **@aliouniang**