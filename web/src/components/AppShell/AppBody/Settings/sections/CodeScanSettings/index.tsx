import { useState, useEffect } from "react";
import {
  Container,
  Title,
  Stack,
  Textarea,
  Switch,
  Group,
  TextInput,
  LoadingOverlay,
} from "@mantine/core";
import {
  useSettings,
  useUpdateSettings,
} from "../../../../../../hooks/use-settings";

export function CodeScanSettings() {
  const { data: settings, isSuccess } = useSettings();
  const updateSettings = useUpdateSettings();

  const [excludeDirs, setExcludeDirs] = useState("");
  const [excludeFiles, setExcludeFiles] = useState("");

  useEffect(() => {
    if (isSuccess && settings) {
      setExcludeDirs(
        settings.code_scan_settings.exclude_directories.join("\n")
      );
      setExcludeFiles(settings.code_scan_settings.exclude_files.join("\n"));
    }
  }, [isSuccess, settings]);

  const handleSyncEnabled = (checked: boolean) => {
    updateSettings({
      code_scan_settings: {
        ...settings!.code_scan_settings,
        sync_enabled: checked,
      },
    });
  };

  const handleGithubAuthChange = (key: string, value: string) => {
    updateSettings({
      github_auth: { ...settings?.github_auth, [key as "token"]: value },
    });
  };

  const handleBlurDirs = () => {
    updateSettings({
      code_scan_settings: {
        ...settings!.code_scan_settings,
        exclude_directories: excludeDirs.split("\n"),
      },
    });
  };

  const handleBlurFiles = () => {
    updateSettings({
      code_scan_settings: {
        ...settings!.code_scan_settings,
        exclude_files: excludeFiles.split("\n"),
      },
    });
  };

  if (!isSuccess) return <LoadingOverlay />;

  return (
    <Container p="xs" fluid>
      <Title size="h3">Code Comment Scanning</Title>
      <Stack gap="md" px="xs">
        <Group grow align="flex-start">
          <Textarea
            autosize
            label="Exclude Directories"
            description="Directories to ignore (one per line)"
            value={excludeDirs}
            onChange={(e) => setExcludeDirs(e.currentTarget.value)}
            onBlur={handleBlurDirs}
          />
          <Textarea
            autosize
            label="Exclude Files"
            description="File patterns to ignore (one per line)"
            value={excludeFiles}
            onChange={(e) => setExcludeFiles(e.currentTarget.value)}
            onBlur={handleBlurFiles}
          />
        </Group>
        <Switch
          label="Sync issues to GitHub"
          checked={settings?.code_scan_settings.sync_enabled}
          onChange={(e) => handleSyncEnabled(e.currentTarget.checked)}
        />
        {settings?.code_scan_settings.sync_enabled && (
          <TextInput
            label="GitHub Token"
            value={settings?.github_auth.token}
            onChange={(e) =>
              handleGithubAuthChange("token", e.currentTarget.value)
            }
          />
        )}
      </Stack>
    </Container>
  );
}
