import { expect, lt, type PanelEditPage } from '@grafana/plugin-e2e';

// @grafana/plugin-e2e >= 3.12.0 clicks the viz picker toggle even when the picker is already open
// (Grafana 12.4), which hides the "All visualizations" tab. Only open the picker when the tab is hidden.
export async function setVisualization(panelEditPage: PanelEditPage, visualization: string) {
  const { selectors, grafanaVersion } = panelEditPage.ctx;
  if (lt(grafanaVersion, '12.4.0')) {
    await panelEditPage.setVisualization(visualization);
    return;
  }

  const { components, constants } = selectors;
  const allVisualizationsTab = panelEditPage.getByGrafanaSelector(components.Tab.title(constants.Tab.title));
  const openVizPickerButton = panelEditPage.getByGrafanaSelector(components.PanelEditor.toggleVizPicker);

  await expect(allVisualizationsTab.or(openVizPickerButton).first()).toBeVisible();
  if (!(await allVisualizationsTab.isVisible())) {
    await openVizPickerButton.click();
  }
  await allVisualizationsTab.click();
  await panelEditPage.getByGrafanaSelector(components.PluginVisualization.item(visualization)).click();
  await expect(
    panelEditPage.getByGrafanaSelector(components.PanelEditor.OptionsPane.header),
    `Could not set visualization to ${visualization}. Ensure the panel is installed.`
  ).toHaveText(visualization);
}
