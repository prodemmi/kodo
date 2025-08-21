import { useState, useEffect, useMemo } from "react";
import {
  Container,
  Title,
  Card,
  Group,
  Text,
  Badge,
  Stack,
  Grid,
  Tabs,
  Timeline,
  Progress,
  Button,
  Alert,
  Loader,
  Center,
  Paper,
  Divider,
  ScrollArea,
  Table,
  ThemeIcon,
  Accordion,
  Select,
  Switch,
  LoadingOverlay,
  Box,
} from "@mantine/core";
import {
  IconHistory,
  IconTrendingUp,
  IconGitBranch,
  IconGitCommit,
  IconRefresh,
  IconCheck,
  IconProgress,
  IconClock,
  IconChevronRight,
  IconInfoCircle,
  IconTrash,
  IconChartLine,
  IconCalendar,
} from "@tabler/icons-react";
import {
  useChanges,
  useCleanupStats,
  useComparison,
  useHistory,
  useRefreshStats,
  useTrends,
} from "../../../../hooks/use-stats";
import { RoleGuard } from "../../../Investor";

// TODO: test todo!
// DONE 2025-08-20 21:11 by prodemmi
export default function History() {
  const [activeTab, setActiveTab] = useState("timeline");
  const [selectedBranch, setSelectedBranch] = useState<string>("all");

  const { data: history, isLoading: isLoadingHistory } = useHistory(true);
  const { data: trends, isLoading: isLoadingTrends } = useTrends(
    activeTab === "trends"
  );
  const { data: changes, isLoading: isLoadingChanges } = useChanges(
    activeTab === "changes"
  );
  const { data: comparison, isLoading: isLoadingComparison } = useComparison(
    activeTab === "comparison"
  );
  const { mutate: refreshStats, isPending: refreshing } = useRefreshStats();
  const { mutate: cleanupStats } = useCleanupStats();

  const loading = useMemo(
    () =>
      isLoadingHistory ||
      isLoadingTrends ||
      isLoadingChanges ||
      isLoadingComparison,
    [isLoadingHistory, isLoadingTrends, isLoadingChanges, isLoadingComparison]
  );

  const getStatusColor = (status: string) => {
    switch (status) {
      case "done":
        return "green";
      case "in-progress":
        return "blue";
      case "todo":
        return "gray";
      default:
        return "gray";
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case "done":
        return IconCheck;
      case "in-progress":
        return IconProgress;
      case "todo":
        return IconClock;
      default:
        return IconClock;
    }
  };

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case "high":
        return "red";
      case "medium":
        return "yellow";
      case "low":
        return "green";
      default:
        return "gray";
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };

  const getFilteredHistory = () => {
    let filtered = history?.history || [];

    if (selectedBranch !== "all") {
      filtered = filtered.filter((h) => h.branch === selectedBranch);
    }

    return filtered.reverse();
  };

  const getUniqueBranches = () => {
    const branches = [
      ...new Set((history?.history || []).map((h) => h.branch)),
    ];
    return branches
      .filter((b) => b !== "unknown")
      .map((branch) => ({ value: branch, label: branch }));
  };

  if (loading) {
    return (
      <LoadingOverlay
        visible={loading}
        zIndex={1000}
        overlayProps={{ radius: "sm", blur: 2 }}
      />
    );
  }

  return (
    <Container size="xl" py="xl">
      <Group justify="space-between" mb="xl">
        <Title order={1}>
          <Group>
            <IconHistory size={32} />
            Board History
          </Group>
        </Title>

        <RoleGuard.Consumer>
          <Group>
            <Button
              leftSection={<IconRefresh size={16} />}
              onClick={() => refreshStats()}
              loading={refreshing}
              variant="light"
            >
              Refresh Stats
            </Button>
            <Button
              leftSection={<IconTrash size={16} />}
              onClick={() => cleanupStats()}
              variant="light"
              color="orange"
            >
              Cleanup Old Stats
            </Button>
          </Group>
        </RoleGuard.Consumer>
      </Group>

      {/* Summary Cards */}
      <Grid mb="xl">
        <Grid.Col span={{ base: 12, md: 3 }}>
          <Card withBorder>
            <Stack align="center">
              <ThemeIcon size="xl" variant="light">
                <IconHistory />
              </ThemeIcon>
              <Text size="xl" fw={700}>
                {history?.count || 0}
              </Text>
              <Text c="dimmed" ta="center">
                Total Snapshots
              </Text>
            </Stack>
          </Card>
        </Grid.Col>

        <Grid.Col span={{ base: 12, md: 3 }}>
          <Card withBorder>
            <Stack align="center">
              <ThemeIcon size="xl" variant="light" color="blue">
                <IconGitBranch />
              </ThemeIcon>
              <Text size="xl" fw={700}>
                {getUniqueBranches().length}
              </Text>
              <Text c="dimmed" ta="center">
                Branches Tracked
              </Text>
            </Stack>
          </Card>
        </Grid.Col>

        {changes && (
          <>
            <Grid.Col span={{ base: 12, md: 3 }}>
              <Card withBorder>
                <Stack align="center">
                  <ThemeIcon size="xl" variant="light" color="green">
                    <IconTrendingUp />
                  </ThemeIcon>
                  <Text size="xl" fw={700}>
                    {changes.summary.added - changes.summary.removed}
                  </Text>
                  <Text c="dimmed" ta="center">
                    Net Change
                  </Text>
                </Stack>
              </Card>
            </Grid.Col>

            <Grid.Col span={{ base: 12, md: 3 }}>
              <Card withBorder>
                <Stack align="center">
                  <ThemeIcon size="xl" variant="light" color="orange">
                    <IconCheck />
                  </ThemeIcon>
                  <Text size="xl" fw={700}>
                    {changes.summary.status_changed}
                  </Text>
                  <Text c="dimmed" ta="center">
                    Status Changes
                  </Text>
                </Stack>
              </Card>
            </Grid.Col>
          </>
        )}
      </Grid>

      <Tabs value={activeTab} onChange={(value) => setActiveTab(value!)}>
        <Tabs.List>
          <Tabs.Tab value="timeline" leftSection={<IconHistory size={16} />}>
            Timeline
          </Tabs.Tab>
          <Tabs.Tab value="trends" leftSection={<IconChartLine size={16} />}>
            Trends
          </Tabs.Tab>
          <Tabs.Tab value="changes" leftSection={<IconTrendingUp size={16} />}>
            Recent Changes
          </Tabs.Tab>
          <Tabs.Tab
            value="comparison"
            leftSection={<IconGitCommit size={16} />}
          >
            Comparison
          </Tabs.Tab>
        </Tabs.List>

        <Tabs.Panel value="timeline" pt="xl">
          <Card withBorder>
            <Select
              mb="lg"
              placeholder="All branches"
              value={selectedBranch}
              onChange={(value) => setSelectedBranch(value || "all")}
              data={[
                { value: "all", label: "All Branches" },
                ...getUniqueBranches(),
              ]}
              leftSection={<IconGitBranch size={16} />}
              w={200}
            />
            <Card.Section p="md">
              <Timeline active={-1} bulletSize={24} lineWidth={2} pr="md">
                {getFilteredHistory().map((snapshot, index) => {
                  const completionRate =
                    snapshot.stats.total > 0
                      ? (snapshot.stats.done / snapshot.stats.total) * 100
                      : 0;

                  return (
                    <Timeline.Item
                      key={`${snapshot.commit}-${index}`}
                      bullet={<IconGitCommit size={12} />}
                      title={
                        <Group>
                          <Badge
                            variant="light"
                            leftSection={<IconGitBranch size={12} />}
                          >
                            {snapshot.branch}
                          </Badge>
                          <Badge variant="outline" c="dimmed">
                            {snapshot.commit_short}
                          </Badge>
                        </Group>
                      }
                    >
                      <Stack gap="xs">
                        <Text size="sm" c="dimmed">
                          {formatDate(snapshot.timestamp)}
                        </Text>
                        {snapshot.commit_message && (
                          <Group gap="2">
                            <Text fw="bold">Commit:</Text>{" "}
                            <Text size="sm">{snapshot.commit_message}</Text>
                          </Group>
                        )}

                        <Grid>
                          <Grid.Col span={6}>
                            <Paper p="xs" withBorder>
                              <Text size="xs" c="dimmed" mb={4}>
                                Progress
                              </Text>
                              <Progress
                                value={completionRate}
                                size="sm"
                                color="green"
                              />
                              <Text size="xs" mt={2}>
                                {completionRate.toFixed(1)}% complete
                              </Text>
                            </Paper>
                          </Grid.Col>

                          <Grid.Col span={6}>
                            <Group gap="xs">
                              <Badge color="gray" size="sm">
                                {snapshot.stats.total} total
                              </Badge>
                              <Badge color="green" size="sm">
                                {snapshot.stats.done} done
                              </Badge>
                              <Badge color="blue" size="sm">
                                {snapshot.stats.in_progress} in progress
                              </Badge>
                              <Badge color="orange" size="sm">
                                {snapshot.stats.todo} todo
                              </Badge>
                            </Group>
                          </Grid.Col>
                        </Grid>

                        {Object.keys(snapshot.stats.by_type).length > 0 && (
                          <Accordion variant="contained">
                            <Accordion.Item value="details">
                              <Accordion.Control>
                                <Text size="sm">View item types</Text>
                              </Accordion.Control>
                              <Accordion.Panel>
                                <Group gap="xs">
                                  {Object.entries(snapshot.stats.by_type).map(
                                    ([type, count]) => (
                                      <Badge
                                        key={type}
                                        variant="outline"
                                        size="sm"
                                      >
                                        {type}: {count}
                                      </Badge>
                                    )
                                  )}
                                </Group>
                              </Accordion.Panel>
                            </Accordion.Item>
                          </Accordion>
                        )}
                      </Stack>
                    </Timeline.Item>
                  );
                })}
              </Timeline>
            </Card.Section>
          </Card>
        </Tabs.Panel>

        <Tabs.Panel value="trends" pt="xl">
          {trends ? (
            <Stack>
              {/* Completion Rate Trend */}
              <Card withBorder>
                <Title order={3} mb="md">
                  Completion Rate Over Time
                </Title>
                <ScrollArea>
                  <Table>
                    <Table.Thead>
                      <Table.Tr>
                        <Table.Th>Commit</Table.Th>
                        <Table.Th>Date</Table.Th>
                        <Table.Th>Completion Rate</Table.Th>
                        <Table.Th>Progress</Table.Th>
                      </Table.Tr>
                    </Table.Thead>
                    <Table.Tbody>
                      {trends.completion_rate.slice(-10).map((entry, index) => (
                        <Table.Tr key={index}>
                          <Table.Td>
                            <Badge variant="outline">{entry.commit}</Badge>
                          </Table.Td>
                          <Table.Td>
                            <Text size="sm">{formatDate(entry.timestamp)}</Text>
                          </Table.Td>
                          <Table.Td>
                            <Text fw={500}>{entry.rate.toFixed(1)}%</Text>
                          </Table.Td>
                          <Table.Td>
                            <Progress value={entry.rate} size="sm" />
                          </Table.Td>
                        </Table.Tr>
                      ))}
                    </Table.Tbody>
                  </Table>
                </ScrollArea>
              </Card>

              {/* Type Trends */}
              <Card withBorder>
                <Title order={3} mb="md">
                  Item Type Trends
                </Title>
                <Accordion>
                  {Object.entries(trends.type_trends).map(([type, data]) => (
                    <Accordion.Item key={type} value={type}>
                      <Accordion.Control>
                        <Group>
                          <Text>{type}</Text>
                          <Badge>{data.length} snapshots</Badge>
                        </Group>
                      </Accordion.Control>
                      <Accordion.Panel>
                        <ScrollArea>
                          <Table>
                            <Table.Thead>
                              <Table.Tr>
                                <Table.Th>Commit</Table.Th>
                                <Table.Th>Date</Table.Th>
                                <Table.Th>Count</Table.Th>
                              </Table.Tr>
                            </Table.Thead>
                            <Table.Tbody>
                              {data.slice(-5).map((entry, index) => (
                                <Table.Tr key={index}>
                                  <Table.Td>
                                    <Badge variant="outline">
                                      {entry.commit}
                                    </Badge>
                                  </Table.Td>
                                  <Table.Td>
                                    <Text size="sm">
                                      {formatDate(entry.timestamp)}
                                    </Text>
                                  </Table.Td>
                                  <Table.Td>
                                    <Badge>{entry.count}</Badge>
                                  </Table.Td>
                                </Table.Tr>
                              ))}
                            </Table.Tbody>
                          </Table>
                        </ScrollArea>
                      </Accordion.Panel>
                    </Accordion.Item>
                  ))}
                </Accordion>
              </Card>
            </Stack>
          ) : (
            <Alert icon={<IconInfoCircle />} title="No trend data available">
              Not enough history snapshots to generate trend analysis. Trends
              require at least 2 snapshots.
            </Alert>
          )}
        </Tabs.Panel>

        <Tabs.Panel value="changes" pt="xl">
          {changes ? (
            <Stack>
              {/* Summary */}
              <Grid>
                <Grid.Col span={4}>
                  <Card withBorder p="md" bg="green.0">
                    <Stack align="center">
                      <ThemeIcon color="green" size="lg">
                        <IconTrendingUp />
                      </ThemeIcon>
                      <Text size="xl" fw={700}>
                        {changes.summary.added}
                      </Text>
                      <Text c="green">Items Added</Text>
                    </Stack>
                  </Card>
                </Grid.Col>

                <Grid.Col span={4}>
                  <Card withBorder p="md" bg="red.0">
                    <Stack align="center">
                      <ThemeIcon color="red" size="lg">
                        <IconTrash />
                      </ThemeIcon>
                      <Text size="xl" fw={700}>
                        {changes.summary.removed}
                      </Text>
                      <Text c="red">Items Removed</Text>
                    </Stack>
                  </Card>
                </Grid.Col>

                <Grid.Col span={4}>
                  <Card withBorder p="md" bg="blue.0">
                    <Stack align="center">
                      <ThemeIcon color="blue" size="lg">
                        <IconCheck />
                      </ThemeIcon>
                      <Text size="xl" fw={700}>
                        {changes.summary.status_changed}
                      </Text>
                      <Text c="blue">Status Changes</Text>
                    </Stack>
                  </Card>
                </Grid.Col>
              </Grid>

              {/* Detailed Changes */}
              <Accordion>
                {changes.added && changes.added.length > 0 && (
                  <Accordion.Item value="added">
                    <Accordion.Control>
                      <Group>
                        <ThemeIcon color="green" size="sm">
                          <IconTrendingUp size={12} />
                        </ThemeIcon>
                        <Text>Added Items ({changes.added.length})</Text>
                      </Group>
                    </Accordion.Control>
                    <Accordion.Panel>
                      <Stack gap="xs">
                        {changes.added.map((item) => (
                          <Paper key={item.id} p="sm" withBorder>
                            <Group justify="space-between">
                              <Group>
                                <Badge variant="light">{item.type}</Badge>
                                <Text>{item.title}</Text>
                              </Group>
                              <Group gap="xs">
                                <Badge
                                  color={getPriorityColor(item.priority)}
                                  size="sm"
                                >
                                  {item.priority}
                                </Badge>
                                <Text size="sm" c="dimmed">
                                  {item.file}:{item.line}
                                </Text>
                              </Group>
                            </Group>
                          </Paper>
                        ))}
                      </Stack>
                    </Accordion.Panel>
                  </Accordion.Item>
                )}

                {changes.removed && changes.removed.length > 0 && (
                  <Accordion.Item value="removed">
                    <Accordion.Control>
                      <Group>
                        <ThemeIcon color="red" size="sm">
                          <IconTrash size={12} />
                        </ThemeIcon>
                        <Text>Removed Items ({changes.removed.length})</Text>
                      </Group>
                    </Accordion.Control>
                    <Accordion.Panel>
                      <Stack gap="xs">
                        {changes.removed.map((item) => (
                          <Paper key={item.id} p="sm" withBorder>
                            <Group justify="space-between">
                              <Group>
                                <Badge variant="light">{item.type}</Badge>
                                <Text>{item.title}</Text>
                              </Group>
                              <Group gap="xs">
                                <Badge
                                  color={getPriorityColor(item.priority)}
                                  size="sm"
                                >
                                  {item.priority}
                                </Badge>
                                <Text size="sm" c="dimmed">
                                  {item.file}:{item.line}
                                </Text>
                              </Group>
                            </Group>
                          </Paper>
                        ))}
                      </Stack>
                    </Accordion.Panel>
                  </Accordion.Item>
                )}

                {changes.status_changed &&
                  changes.status_changed.length > 0 && (
                    <Accordion.Item value="status-changed">
                      <Accordion.Control>
                        <Group>
                          <ThemeIcon color="blue" size="sm">
                            <IconCheck size={12} />
                          </ThemeIcon>
                          <Text>
                            Status Changes ({changes.status_changed.length})
                          </Text>
                        </Group>
                      </Accordion.Control>
                      <Accordion.Panel>
                        <Stack gap="xs">
                          {changes.status_changed.map((change, index) => (
                            <Paper key={index} p="sm" withBorder>
                              <Group justify="space-between">
                                <Group>
                                  <Badge variant="light">
                                    {change.item.type}
                                  </Badge>
                                  <Text>{change.item.title}</Text>
                                </Group>
                                <Group gap="xs">
                                  <Badge
                                    color={getStatusColor(
                                      change.old_status || ""
                                    )}
                                    size="sm"
                                  >
                                    {change.old_status}
                                  </Badge>
                                  <IconChevronRight size={12} />
                                  <Badge
                                    color={getStatusColor(
                                      change.new_status || ""
                                    )}
                                    size="sm"
                                  >
                                    {change.new_status}
                                  </Badge>
                                  <Text size="sm" c="dimmed">
                                    {change.item.file}:{change.item.line}
                                  </Text>
                                </Group>
                              </Group>
                            </Paper>
                          ))}
                        </Stack>
                      </Accordion.Panel>
                    </Accordion.Item>
                  )}
              </Accordion>
            </Stack>
          ) : (
            <Alert icon={<IconInfoCircle />} title="No recent changes">
              No recent changes detected. This requires at least 2 history
              snapshots to compare.
            </Alert>
          )}
        </Tabs.Panel>

        <Tabs.Panel value="comparison" pt="xl">
          {comparison ? (
            <Stack>
              <Grid>
                <Grid.Col span={6}>
                  <Card withBorder>
                    <Stack>
                      <Group>
                        <ThemeIcon color="blue">
                          <IconGitCommit />
                        </ThemeIcon>
                        <div>
                          <Text fw={500}>Current Commit</Text>
                          <Text size="sm" c="dimmed">
                            {comparison.current.commit}
                          </Text>
                        </div>
                      </Group>

                      <Divider />

                      <Group justify="space-around">
                        <Stack align="center">
                          <Text size="xl" fw={700}>
                            {comparison.current.stats.total}
                          </Text>
                          <Text size="sm" c="dimmed">
                            Total
                          </Text>
                        </Stack>
                        <Stack align="center">
                          <Text size="xl" fw={700} c="green">
                            {comparison.current.stats.done}
                          </Text>
                          <Text size="sm" c="dimmed">
                            Done
                          </Text>
                        </Stack>
                        <Stack align="center">
                          <Text size="xl" fw={700} c="blue">
                            {comparison.current.stats.in_progress}
                          </Text>
                          <Text size="sm" c="dimmed">
                            In Progress
                          </Text>
                        </Stack>
                        <Stack align="center">
                          <Text size="xl" fw={700} c="orange">
                            {comparison.current.stats.todo}
                          </Text>
                          <Text size="sm" c="dimmed">
                            Todo
                          </Text>
                        </Stack>
                      </Group>

                      <Progress
                        value={
                          comparison.current.stats.total > 0
                            ? (comparison.current.stats.done /
                                comparison.current.stats.total) *
                              100
                            : 0
                        }
                        color="green"
                        size="md"
                      />
                    </Stack>
                  </Card>
                </Grid.Col>

                <Grid.Col span={6}>
                  <Card withBorder>
                    <Stack>
                      <Group>
                        <ThemeIcon color="gray">
                          <IconGitCommit />
                        </ThemeIcon>
                        <div>
                          <Text fw={500}>Previous Commit</Text>
                          <Text size="sm" c="dimmed">
                            {comparison.previous.commit}
                          </Text>
                        </div>
                      </Group>

                      <Divider />

                      <Group justify="space-around">
                        <Stack align="center">
                          <Text size="xl" fw={700}>
                            {comparison.previous.stats.total}
                          </Text>
                          <Text size="sm" c="dimmed">
                            Total
                          </Text>
                        </Stack>
                        <Stack align="center">
                          <Text size="xl" fw={700} c="green">
                            {comparison.previous.stats.done}
                          </Text>
                          <Text size="sm" c="dimmed">
                            Done
                          </Text>
                        </Stack>
                        <Stack align="center">
                          <Text size="xl" fw={700} c="blue">
                            {comparison.previous.stats.in_progress}
                          </Text>
                          <Text size="sm" c="dimmed">
                            In Progress
                          </Text>
                        </Stack>
                        <Stack align="center">
                          <Text size="xl" fw={700} c="orange">
                            {comparison.previous.stats.todo}
                          </Text>
                          <Text size="sm" c="dimmed">
                            Todo
                          </Text>
                        </Stack>
                      </Group>

                      <Progress
                        value={
                          comparison.previous.stats.total > 0
                            ? (comparison.previous.stats.done /
                                comparison.previous.stats.total) *
                              100
                            : 0
                        }
                        color="green"
                        size="md"
                      />
                    </Stack>
                  </Card>
                </Grid.Col>
              </Grid>

              {/* Changes Summary */}
              <Card withBorder>
                <Title order={3} mb="md">
                  Changes Summary
                </Title>
                <Group grow>
                  <Paper
                    p="md"
                    bg={
                      comparison.changes.total > 0
                        ? "green.0"
                        : comparison.changes.total < 0
                        ? "red.0"
                        : "gray.0"
                    }
                  >
                    <Stack align="center">
                      <Text size="xl" fw={700}>
                        {comparison.changes.total > 0
                          ? `+${comparison.changes.total}`
                          : comparison.changes.total}
                      </Text>
                      <Text c="dimmed">Total Items</Text>
                    </Stack>
                  </Paper>

                  <Paper
                    p="md"
                    bg={
                      comparison.changes.done > 0
                        ? "green.0"
                        : comparison.changes.done < 0
                        ? "red.0"
                        : "gray.0"
                    }
                  >
                    <Stack align="center">
                      <Text size="xl" fw={700}>
                        {comparison.changes.done > 0
                          ? `+${comparison.changes.done}`
                          : comparison.changes.done}
                      </Text>
                      <Text c="dimmed">Done Items</Text>
                    </Stack>
                  </Paper>

                  <Paper
                    p="md"
                    bg={
                      comparison.changes.in_progress > 0
                        ? "blue.0"
                        : comparison.changes.in_progress < 0
                        ? "red.0"
                        : "gray.0"
                    }
                  >
                    <Stack align="center">
                      <Text size="xl" fw={700}>
                        {comparison.changes.in_progress > 0
                          ? `+${comparison.changes.in_progress}`
                          : comparison.changes.in_progress}
                      </Text>
                      <Text c="dimmed">In Progress</Text>
                    </Stack>
                  </Paper>

                  <Paper
                    p="md"
                    bg={
                      comparison.changes.todo > 0
                        ? "orange.0"
                        : comparison.changes.todo < 0
                        ? "green.0"
                        : "gray.0"
                    }
                  >
                    <Stack align="center">
                      <Text size="xl" fw={700}>
                        {comparison.changes.todo > 0
                          ? `+${comparison.changes.todo}`
                          : comparison.changes.todo}
                      </Text>
                      <Text c="dimmed">Todo Items</Text>
                    </Stack>
                  </Paper>
                </Group>
              </Card>
            </Stack>
          ) : (
            <Alert icon={<IconInfoCircle />} title="No comparison data">
              Not enough history to compare commits. Comparison requires at
              least 2 snapshots.
            </Alert>
          )}
        </Tabs.Panel>
      </Tabs>

      {/* Footer Stats */}
      {history?.count && (
        <Card withBorder mt="xl">
          <Group justify="space-between">
            <Group>
              <IconCalendar size={16} />
              <Text size="sm" c="dimmed">
                First snapshot: {formatDate(history.history[0].timestamp)}
              </Text>
            </Group>
            <Group>
              <IconHistory size={16} />
              <Text size="sm" c="dimmed">
                Last updated:{" "}
                {formatDate(history.history[history.count - 1].timestamp)}
              </Text>
            </Group>
          </Group>
        </Card>
      )}
    </Container>
  );
}
