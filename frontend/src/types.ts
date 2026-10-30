export type RoomStatus = 'lobby' | 'loading' | 'starting' | 'round_active' | 'round_reveal' | 'game_over';

export interface ParameterItem {
  label: string;
  value: string;
}

export interface Player {
  id: string;
  nickname: string;
  avatar: string;
  is_ready: boolean;
  score: number;
  has_voted: boolean;
  has_guess: boolean;
}

export interface RoundListing {
  id: string;
  category: string;
  title: string;
  description: string;
  image_path: string;
  region: string;
  parameters?: ParameterItem[];
  options?: number[];
  actual_price?: number;
}

export interface PlayerGuessResult {
  player_id: string;
  nickname: string;
  avatar: string;
  guess_value: number;
  accuracy_pct: number;
  score_points: number;
  is_bullseye: boolean;
  time_taken_ms: number;
}

export interface Room {
  id: string;
  code: string;
  host_id: string;
  is_public: boolean;
  category: string;
  game_mode: 'exact' | 'multiple_choice' | 'mix';
  question_count: number;
  time_limit_sec: number;
  max_players: number;
  status: RoomStatus;
  players: Record<string, Player>;
  votes_count: Record<string, number>;
  votes_mode: Record<string, string>;
  votes_category: Record<string, string>;
  current_round: number;
  round_started_at: string;
  round_guesses: Record<string, PlayerGuessResult>;
  current_listing?: RoundListing;
}

export interface PublicRoomItem {
  code: string;
  category: string;
  game_mode: string;
  player_count: number;
  max_players: number;
  question_count: number;
  status: string;
}

