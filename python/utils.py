class WarehouseSize:
    def __init__(self, size: str, min_size: str, max_size: str):
        self.available_size: list[str] = [
            "xsmall",
            "small",
            "medium",
            "large",
            "xlarge",
            "xxlarge",
            "xxxlarge",
            "xxxxlarge",
        ]
        self.size: str = self.size_caster(size)
        self.min_size: str = self.size_caster(min_size)
        self.max_size: str = self.size_caster(max_size)
        self.min_index: int = self.available_size.index(self.min_size)
        self.max_index: int = self.available_size.index(self.max_size)

        if self.min_index > self.max_index:
            raise ValueError(
                "Min warehouse size cannot be larger than max warehouse size"
            )

    def upsize(self):
        index: int = self.available_size.index(self.size)
        # if current_size > max_size, go to nearest max_size
        # this is to ensure the wh stay max capped size
        if index > self.max_index:
            self.size = self.available_size[self.max_index]
            return True
        elif index < self.max_index:
            self.size = self.available_size[index + 1]
            return True
        return False

    def downsize(self):
        index: int = self.available_size.index(self.size)
        # if current_size < min_size, go to nearest min_size
        # this is to ensure the wh stay min floored size
        if index < self.min_index:
            self.size = self.available_size[self.min_index]
            return True
        elif index > self.min_index:
            self.size = self.available_size[index - 1]
            return True
        return False

    def size_caster(self, size: str):
        match size:
            case "X-Small", "xsmall":
                size = "xsmall"
            case "Small" | "small":
                size = "small"
            case "Medium" | "medium":
                size = "medium"
            case "Large" | "large":
                size = "large"
            case "X-Large" | "xlarge":
                size = "xlarge"
            case "2X-Large" | "xxlarge":
                size = "xxlarge"
            case "3X-Large" | "xxxlarge":
                size = "xxxlarge"
            case "4X-Large" | "xxxxlarge":
                size = "xxxxlarge"
            case _:
                size = "xsmall"
        return size
