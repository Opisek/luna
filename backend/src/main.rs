use axum::{
    routing::get,
    Router,
};

#[tokio::main]
async fn main() {
    let endpoints = Router::new().route("/calendars", get(|| async { "Hello, World!" }));

    let app = Router::new().nest("/api", endpoints);
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}