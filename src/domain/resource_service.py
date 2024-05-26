from domain.resource import Resource

def createResource(key: str, resources: dict) -> tuple:
    resource = __lastResource(key, resources)
    if resource == None:
        resource = Resource()
        resources[key] = [resource]
    elif resource.isLast():
        resource = Resource()
        resources[key] = [*resources[key], resource]
    return resources, resource
    
def updateLastResourceByKey(key: str, resources: dict, resource: Resource) -> dict:
    if key not in resources:
        return resources
    index = len(resources[key]) - 1
    resources[key][index] = resource
    return resources

def __lastResource(key: str, resources: dict) -> Resource:
    if key not in resources:
        return None
    
    index = len(resources[key]) - 1
    return resources[key][index]

