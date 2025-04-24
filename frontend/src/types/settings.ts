export enum UserSettingKeys {
  DebugMode = "debug_mode",
  DisplayWeekNumbers = "display_week_numbers",
  FirstDayOfWeek = "first_day_of_week",
  ThemeLight = "theme_light",
  ThemeDark = "theme_dark",
  ThemeSynchronize = "theme_sync",
  FontText = "font_text",
  FontTime = "font_time",
  DisplayAllDayEventsFilled = "display_all_day_events_filled",
  DisplayNonAllDayEventsFilled = "display_non_all_day_events_filled",
  DisplaySmallCalendar = "display_small_calendar",
  DynamicCalendarRows = "dynamic_calendar_rows",
  DynamicSmallCalendarRows = "dynamic_small_calendar_rows",
  DisplayRoundedCorners = "display_rounded_corners",
  UiScaling = "ui_scaling",
  AnimateCalendarSwipe = "animate_calendar_swipe",
  AnimateSmallCalendarSwipe = "animate_small_calendar_swipe",
  AnimateMonthSelectionSwipe = "animate_month_selection_swipe",
  AppearenceFrostedGlass = "appearance_frosted_glass",
}

export enum GlobalSettingKeys {
  RegistrationEnabled = "registration_enabled",
  LoggingVerbosity = "logging_verbosity",
  UseCdnFonts = "use_cdn_fonts",
}

export type UserSettings = {
  [UserSettingKeys.DebugMode]: boolean;
  [UserSettingKeys.DisplayWeekNumbers]: boolean;
  [UserSettingKeys.FirstDayOfWeek]: number;
  [UserSettingKeys.ThemeLight]: string;
  [UserSettingKeys.ThemeDark]: string;
  [UserSettingKeys.FontText]: string;
  [UserSettingKeys.FontTime]: string;
  [UserSettingKeys.ThemeSynchronize]: boolean;
  [UserSettingKeys.DisplayAllDayEventsFilled]: boolean;
  [UserSettingKeys.DisplayNonAllDayEventsFilled]: boolean;
  [UserSettingKeys.DisplaySmallCalendar]: boolean;
  [UserSettingKeys.DynamicCalendarRows]: boolean;
  [UserSettingKeys.DynamicSmallCalendarRows]: boolean;
  [UserSettingKeys.DisplayRoundedCorners]: boolean;
  [UserSettingKeys.UiScaling]: number;
  [UserSettingKeys.AnimateCalendarSwipe]: boolean;
  [UserSettingKeys.AnimateSmallCalendarSwipe]: boolean;
  [UserSettingKeys.AnimateMonthSelectionSwipe]: boolean;
  [UserSettingKeys.AppearenceFrostedGlass]: boolean;
};

export type GlobalSettings = {
  [GlobalSettingKeys.RegistrationEnabled]: boolean;
  [GlobalSettingKeys.LoggingVerbosity]: number;
  [GlobalSettingKeys.UseCdnFonts]: boolean;
};

export type UserData = {
  id: string;
  username: string;
  email: string;
  searchable: boolean;
  profile_picture: string;
  admin: boolean;
}