import React, { useContext, useMemo, useCallback } from "react";
import { useParams } from "react-router-dom";
import {
  DetailsList,
  IColumn,
  SelectionMode,
  ActionButton,
} from "@fluentui/react";
import { FormattedMessage, Context } from "@oursky/react-messageformat";
import ScreenContent from "../../ScreenContent";
import ScreenTitle from "../../ScreenTitle";
import ScreenDescription from "../../ScreenDescription";
import Widget from "../../Widget";
import WidgetTitle from "../../WidgetTitle";
import ShowLoading from "../../ShowLoading";
import ShowError from "../../ShowError";
import {
  useAppAndSecretConfigQuery,
  AppAndSecretConfigQueryResult,
} from "./query/appAndSecretConfigQuery";
import { formatDatetime } from "../../util/formatDatetime";
import { useSystemConfig } from "../../context/SystemConfigContext";
import { downloadStringAsFile } from "../../util/download";
import styles from "./AdminAPIConfigurationScreen.module.scss";

interface AdminAPIConfigurationScreenContentProps {
  queryResult: AppAndSecretConfigQueryResult;
}

interface Item {
  keyID: string;
  createdAt: string | null;
  publicKeyPEM: string;
  privateKeyPEM?: string | null;
}

const AdminAPIConfigurationScreenContent: React.FC<AdminAPIConfigurationScreenContentProps> =
  function AdminAPIConfigurationScreenContent(props) {
    const { locale, renderToString } = useContext(Context);
    const { themes } = useSystemConfig();

    const adminAPISecrets = useMemo(() => {
      return props.queryResult.secretConfig?.adminAPISecrets ?? [];
    }, [props.queryResult.secretConfig?.adminAPISecrets]);

    const items: Item[] = useMemo(() => {
      const items = [];
      for (const adminAPISecret of adminAPISecrets) {
        items.push({
          keyID: adminAPISecret.keyID,
          createdAt: formatDatetime(locale, adminAPISecret.createdAt),
          publicKeyPEM: adminAPISecret.publicKeyPEM,
          privateKeyPEM: adminAPISecret.privateKeyPEM,
        });
      }
      return items;
    }, [locale, adminAPISecrets]);

    const downloadItem = useCallback(
      (keyID: string) => {
        const item = items.find((a) => a.keyID === keyID);
        if (item == null) {
          return;
        }
        if (item.privateKeyPEM != null) {
          downloadStringAsFile({
            content: item.privateKeyPEM,
            mimeType: "application/x-pem-file",
            filename: `${item.keyID}.pem`,
          });
        }
        // TODO: reauthenticate
      },
      [items]
    );

    const actionColumnOnRender = useCallback(
      (item?: Item) => {
        return (
          <ActionButton
            className={styles.actionButton}
            theme={themes.actionButton}
            onClick={(e: React.MouseEvent<unknown>) => {
              e.preventDefault();
              e.stopPropagation();
              if (item != null) {
                downloadItem(item.keyID);
              }
            }}
          >
            <FormattedMessage id="download" />
          </ActionButton>
        );
      },
      [downloadItem, themes.actionButton]
    );

    const columns: IColumn[] = useMemo(() => {
      return [
        {
          key: "keyID",
          fieldName: "keyID",
          name: renderToString("AdminAPIConfigurationScreen.column.key-id"),
          minWidth: 150,
        },
        {
          key: "createdAt",
          fieldName: "createdAt",
          name: renderToString("AdminAPIConfigurationScreen.column.created-at"),
          minWidth: 150,
        },
        {
          key: "action",
          name: renderToString("action"),
          minWidth: 150,
          onRender: actionColumnOnRender,
        },
      ];
    }, [renderToString, actionColumnOnRender]);

    return (
      <ScreenContent className={styles.root}>
        <ScreenTitle>
          <FormattedMessage id="AdminAPIConfigurationScreen.title" />
        </ScreenTitle>
        <ScreenDescription className={styles.widget}>
          <FormattedMessage id="AdminAPIConfigurationScreen.description" />
        </ScreenDescription>
        <Widget className={styles.widget}>
          <WidgetTitle>
            <FormattedMessage id="AdminAPIConfigurationScreen.keys.title" />
          </WidgetTitle>
          <DetailsList
            items={items}
            columns={columns}
            selectionMode={SelectionMode.none}
          />
        </Widget>
      </ScreenContent>
    );
  };

const AdminAPIConfigurationScreen: React.FC =
  function AdminAPIConfigurationScreen() {
    const { appID } = useParams();
    const queryResult = useAppAndSecretConfigQuery(appID);

    if (queryResult.loading) {
      return <ShowLoading />;
    }

    if (queryResult.error) {
      return (
        <ShowError error={queryResult.error} onRetry={queryResult.refetch} />
      );
    }

    return <AdminAPIConfigurationScreenContent queryResult={queryResult} />;
  };

export default AdminAPIConfigurationScreen;
