import { Item } from "./item";

export interface BranchSnapshot {
  branch: string;
  commit: string;
  commit_short: string;
  commit_message: string;
  timestamp: string;
  history: {
    total: number;
    todo: number;
    in_progress: number;
    done: number;
    by_type: Record<string, number>;
    by_priority: Record<string, number>;
    items: Item[];
  };
}

export interface ItemChange {
  item: Item;
  old_status?: string;
  new_status?: string;
}

export interface TrendData {
  timeline: Array<{
    timestamp: string;
    commit: string;
    branch: string;
    total: number;
    todo: number;
    in_progress: number;
    done: number;
  }>;
  completion_rate: Array<{
    timestamp: string;
    commit: string;
    rate: number;
  }>;
  type_trends: Record<
    string,
    Array<{
      timestamp: string;
      commit: string;
      count: number;
    }>
  >;
}

export interface RecentChanges {
  added?: Item[];
  removed?: Item[];
  status_changed?: ItemChange[];
  summary: {
    added: number;
    removed: number;
    status_changed: number;
  };
}

export interface History {
  count: number;
  history: BranchSnapshot[];
}

export interface Changes {
  done: number;
  in_progress: number;
  todo: number;
  total: number;
}

export interface Compare {
  changes: Changes;
  current: Omit<BranchSnapshot, "commit_short", "commit_message">;
  previous: Omit<BranchSnapshot, "commit_short", "commit_message">;
}
