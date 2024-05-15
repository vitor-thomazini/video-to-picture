package domain

import "image"

func Draw(frame image.Image, resource Resource) Resource {
	var drawer Drawer
	style := Style{
		MarginX:  0,
		MarginY:  0,
		PaddingX: 0,
		PaddingY: 0,
	}
	if resource.IsFirstPictureToFirstRow() {
		drawer = NewDrawerStartPoint(frame, style)
	} else if resource.IsFirstPictureToAnyRow() {
		drawer = NewDrawerStartPointFromRect(*resource.Background(), frame, style)
	} else {
		drawer = NewDrawerMiddlePoint(*resource.Background(), frame, style)
	}

	resource.UpdateBackground(drawer.CalculatePanel())
	drawer.DrawIn(resource.Image(), resource.Background())
	resource.IncPictureCounter()
	return resource
}

func DrawAndUpdateResources(image image.Image, mapping map[string][]Resource, key string) map[string][]Resource {
	bkgList, bkg := AppendLastResourceToResourceMap(mapping, key)
	bkg = Draw(image, bkg)
	return AppendResourceMap(bkgList, key, bkg)
}
