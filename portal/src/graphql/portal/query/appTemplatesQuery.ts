import { useMemo } from "react";
import { gql, QueryResult, useQuery } from "@apollo/client";

import { client } from "../apollo";
import {
  AppTemplatesQuery,
  AppTemplatesQueryVariables,
} from "./__generated__/AppTemplatesQuery";
import { renderPath } from "../../../templates";
import {
  Resource,
  ResourceDefinition,
  ResourceSpecifier,
  LanguageTag,
  decodeForText,
  binary,
} from "../../../util/resource";

export const appTemplatesQuery = gql`
  query AppTemplatesQuery($id: ID!, $paths: [String!]!) {
    node(id: $id) {
      __typename
      ... on App {
        id
        resources(paths: $paths) {
          path
          languageTag
          data
          effectiveData
        }
      }
    }
  }
`;

export interface AppTemplatesQueryResult
  extends Pick<
    QueryResult<AppTemplatesQuery, AppTemplatesQueryVariables>,
    "loading" | "error" | "refetch"
  > {
  resources: Record<string, Resource>;
}

export function useAppTemplatesQuery(
  appID: string,
  locales: LanguageTag[],
  ...resourceDefs: ResourceDefinition[]
): AppTemplatesQueryResult {
  const specifiers = useMemo<ResourceSpecifier[]>(() => {
    const output: ResourceSpecifier[] = [];
    for (const locale of locales) {
      for (const resourceDef of resourceDefs) {
        output.push({
          locale,
          def: resourceDef,
          path: renderPath(resourceDef.resourcePath, { locale }),
        });
      }
    }
    return output;
  }, [locales, resourceDefs]);

  const paths = useMemo(() => specifiers.map((specifier) => specifier.path), [
    specifiers,
  ]);

  const { data, loading, error, refetch } = useQuery<
    AppTemplatesQuery,
    AppTemplatesQueryVariables
  >(appTemplatesQuery, {
    client,
    variables: {
      id: appID,
      paths,
    },
  });

  // eslint-disable-next-line complexity
  const resources = useMemo(() => {
    const appNode = data?.node?.__typename === "App" ? data.node : null;
    const resources: Record<string, Resource> = {};

    for (const specifier of specifiers) {
      let found = false;

      for (const resource of appNode?.resources ?? []) {
        if (specifier.path === resource.path) {
          found = true;
          let value = "";
          let transform: (a: string) => string;
          switch (specifier.def.type) {
            case "text":
              transform = decodeForText;
              break;
            case "binary":
              transform = binary;
              break;
            default:
              throw new Error(
                "unexpected resource type: " + String(specifier.def.type)
              );
          }

          if (resource.data != null) {
            value = transform(resource.data);
          } else if (
            resource.effectiveData != null &&
            specifier.def.usesEffectiveDataAsFallbackValue
          ) {
            value = transform(resource.effectiveData);
          }
          resources[specifier.path] = {
            ...specifier,
            value,
          };
          break;
        }
      }

      if (!found) {
        resources[specifier.path] = {
          ...specifier,
          value: "",
        };
      }
    }

    return resources;
  }, [data, specifiers]);

  return { resources, loading, error, refetch };
}
