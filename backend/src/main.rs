use axum::{
  Json,
  Router,
  routing::get,
};
use serde::Serialize;
use minicaldav::{self};
use url::Url;

#[derive(Serialize)]
struct Calendar {
  pub url: Url,
  pub name: String,
  pub color: Option<String>,
}
impl From<&minicaldav::Calendar> for Calendar {
  fn from(calendar: &minicaldav::Calendar) -> Self {
    Calendar {
      url: calendar.url().clone(),
      name: calendar.name().clone(),
      color: calendar.color().cloned(),
    }
  }
}

#[derive(Serialize)]
struct Event {
  name: String,
}

async fn get_calendars() -> Json<Vec<Calendar>> {
  let agent = ureq::Agent::new();
  let url = url::Url::parse("").unwrap();

  let username = "";
  let password = "";
  let credentials = minicaldav::Credentials::Basic(username.into(), password.into());

  let calendars_raw = minicaldav::get_calendars(agent.clone(), &credentials, &url).unwrap();
  let calendars: Vec<Calendar> = calendars_raw.iter().map(|c| Calendar::from(c)).collect();

  Json(calendars)
}

async fn get_events() -> Json<Vec<Event>> {
  let events = vec![
    Event { name: "Event2".to_string() },
    Event { name: "Event1".to_string() },
  ];
  Json(events)
}

#[tokio::main]
async fn main() {
  let endpoints = Router::new().route("/calendars", get(get_calendars));

  let app = Router::new().nest("/api", endpoints);
  let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
  axum::serve(listener, app).await.unwrap();
}